package engines

import (
	"database/sql"
	"fmt"

	"github.com/phassans/frolleague/clients/linkedin"
	"github.com/phassans/frolleague/clients/phantom"
	"github.com/phassans/frolleague/common"
	"github.com/rs/zerolog"
)

type (
	linkedInEngine struct {
		client   linkedin.Client
		sql      *sql.DB
		logger   zerolog.Logger
		pClient  phantom.Client
		dbEngine DatabaseEngine
	}

	LinkedInEngine interface {
		// LinkedInMethods
		GetAccessToken(linkedin.AuthCode) (linkedin.AccessTokenResponse, error)
		LogIn(linkedin.AuthCode, linkedin.AccessToken) (linkedin.MeResponse, error)
		GetMe(linkedin.UserID) (linkedin.MeResponse, bool, error)

		// LinkedInCrawl Methods
		CrawlUserProfile(LinkedInUserID) (phantom.Profile, error)
		SaveUserProfile(LinkedInUserID, phantom.Profile) error
		UpdateUserWithLinkedInURL(UserID LinkedInUserID, url LinkedInURL) error
	}
)

func NewLinkedInEngine(client linkedin.Client, psql *sql.DB, logger zerolog.Logger, pClient phantom.Client, dbEngine DatabaseEngine) LinkedInEngine {
	return &linkedInEngine{client, psql, logger, pClient, dbEngine}
}

func (l *linkedInEngine) GetAccessToken(authCode linkedin.AuthCode) (linkedin.AccessTokenResponse, error) {
	return l.client.GetAccessToken(authCode)
}

func (l *linkedInEngine) LogIn(authCode linkedin.AuthCode, token linkedin.AccessToken) (linkedin.MeResponse, error) {
	// get access token
	if token == "" {
		resp, err := l.client.GetAccessToken(authCode)
		if err != nil {
			return linkedin.MeResponse{}, err
		}
		token = resp.AccessToken
	}

	fmt.Printf("accessToken: %s", token)

	// get me
	meResp, err := l.client.GetMe(token)
	if err != nil {
		return linkedin.MeResponse{}, err
	}

	// save user token
	if err := l.dbEngine.SaveToken(LinkedInUserID(meResp.ID), AccessToken(token)); err != nil {
		return linkedin.MeResponse{}, err
	}

	// save user
	if err := l.dbEngine.SaveUser(LinkedInUserID(meResp.ID), FirstName(meResp.FirstName.Localized.EnUS), LastName(meResp.LastName.Localized.EnUS), LinkedInImage(meResp.ProfilePicture.DisplayImage)); err != nil {
		switch err.(type) {
		case common.DuplicateLinkedInUser:
			fmt.Printf("duplicate user!")
			return meResp, nil
		default:
			return meResp, err
		}
	}

	return meResp, nil
}

func (l *linkedInEngine) GetMe(userID linkedin.UserID) (linkedin.MeResponse, bool, error) {
	token, err := l.dbEngine.GetTokenByUserID(LinkedInUserID(userID))
	if err != nil {
		return linkedin.MeResponse{}, false, nil
	}

	if token != "" {
		meResp, err := l.client.GetMe(linkedin.AccessToken(token))
		if err != nil {
			return linkedin.MeResponse{}, false, err
		}
		if meResp.ID == string(userID) {
			return meResp, true, nil
		}
	}
	return linkedin.MeResponse{}, false, nil
}

func (l *linkedInEngine) CrawlUserProfile(UserID LinkedInUserID) (phantom.Profile, error) {
	user, err := l.dbEngine.GetUserByID(UserID)
	if err != nil {
		return phantom.Profile{}, nil
	}

	// get userProfile
	profile, err := l.pClient.GetUserProfile(string(user.LinkedInURL), false)
	if err != nil {
		return phantom.Profile{}, err
	}

	if err := l.SaveUserProfile(UserID, profile); err != nil {
		return profile, err
	}

	return profile, nil
}

func (l *linkedInEngine) SaveUserProfile(userID LinkedInUserID, profile phantom.Profile) error {
	if err := l.addUserToSchools(profile, UserID(userID)); err != nil {
		return err
	}

	if err := l.addUserToCompanies(profile, UserID(userID)); err != nil {
		return err
	}
	return nil
}

func (l *linkedInEngine) addUserToSchools(profile phantom.Profile, userID UserID) error {
	for _, school := range profile.Schools {
		schoolID, err := l.dbEngine.AddSchoolIfNotPresent(SchoolName(school.SchoolName), Degree(school.Degree), FieldOfStudy(school.FieldOfStudy))
		if err != nil {
			return err
		}

		if err := l.dbEngine.AddUserToSchool(userID, schoolID, FromYear(school.FromYear), ToYear(school.ToYear)); err != nil {
			return err
		}
	}
	return nil
}

func (l *linkedInEngine) addUserToCompanies(profile phantom.Profile, userID UserID) error {
	for _, company := range profile.Companies {
		companyID, err := l.dbEngine.AddCompanyIfNotPresent(CompanyName(company.CompanyName), Location(company.Location))
		if err != nil {
			return err
		}

		if err := l.dbEngine.AddUserToCompany(userID, companyID, Title(company.Title), FromYear(company.FromYear), ToYear(company.ToYear)); err != nil {
			return err
		}

		// update user preferences
		groups, err := l.dbEngine.AddGroupsToUser(userID)
		if err != nil {
			return err
		}
		l.logger.Info().Msgf("groups: %v", groups)

	}

	return nil
}

func (l *linkedInEngine) UpdateUserWithLinkedInURL(UserID LinkedInUserID, url LinkedInURL) error {
	return l.dbEngine.UpdateUserWithLinkedInURL(UserID, url)
}
