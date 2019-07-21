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
		client           linkedin.Client
		sql              *sql.DB
		logger           zerolog.Logger
		pClient          phantom.Client
		dbEngine         DatabaseEngine
		rocketChatEngine RocketChatEngine
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

		// Adds
		UpdateUserSchools(UserID, []phantom.School) error
		UpdateUserCompanies(UserID, []phantom.Company) error
	}
)

func NewLinkedInEngine(
	client linkedin.Client,
	psql *sql.DB,
	logger zerolog.Logger,
	pClient phantom.Client,
	dbEngine DatabaseEngine,
	rocketChatEngine RocketChatEngine,
) LinkedInEngine {
	return &linkedInEngine{client, psql, logger, pClient, dbEngine, rocketChatEngine}
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
		return linkedin.MeResponse{}, false, err
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
	profile, err := l.pClient.GetUserProfile(string(user.LinkedInURL), true)
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
	return l.UpdateUserSchools(userID, profile.Schools)
}

func (l *linkedInEngine) addUserToCompanies(profile phantom.Profile, userID UserID) error {
	return l.UpdateUserCompanies(userID, profile.Companies)
}

func (l *linkedInEngine) UpdateUserWithLinkedInURL(UserID LinkedInUserID, url LinkedInURL) error {
	_, err := l.dbEngine.GetUserByID(UserID)
	if err != nil {
		return err
	}
	return l.dbEngine.UpdateUserWithLinkedInURL(UserID, url)
}

func (l *linkedInEngine) UpdateUserCompanies(userID UserID, companies []phantom.Company) error {
	// clear user-company relationship
	if err := l.dbEngine.UpdateUserStatusForAllCompanies(userID); err != nil {
		return err
	}

	for _, company := range companies {
		companyID, err := l.dbEngine.AddCompanyIfNotPresent(CompanyName(company.CompanyName), Location(company.Location))
		if err != nil {
			return err
		}

		if err := l.dbEngine.AddUserToCompany(userID, companyID, Title(company.Title), FromYear(company.FromYear), ToYear(company.ToYear)); err != nil {
			return err
		}

	}

	// update user preferences
	if err := l.UpdateUserPreferences(userID); err != nil {
		return err
	}

	return nil
}

func (l *linkedInEngine) UpdateUserSchools(userID UserID, schools []phantom.School) error {
	// clear all user-school relationship
	if err := l.dbEngine.UpdateUserStatusForAllSchools(userID); err != nil {
		return err
	}

	for _, school := range schools {
		schoolID, err := l.dbEngine.AddSchoolIfNotPresent(SchoolName(school.SchoolName), Degree(school.Degree), FieldOfStudy(school.FieldOfStudy))
		if err != nil {
			return err
		}

		if err := l.dbEngine.AddUserToSchool(userID, schoolID, FromYear(school.FromYear), ToYear(school.ToYear)); err != nil {
			return err
		}
	}

	// update user preferences
	if err := l.UpdateUserPreferences(userID); err != nil {
		return err
	}

	return nil
}

func (l *linkedInEngine) UpdateUserPreferences(userID UserID) error {
	// update user preferences
	addedGroups, err := l.dbEngine.AddGroupsToUser(userID)
	if err != nil {
		return err
	}
	l.logger.Info().Msgf("added groups: %v", addedGroups)

	if len(addedGroups) > 0 {
		// add rocket chat
		if err := l.rocketChatEngine.AddUserToGroups(userID, addedGroups); err != nil {
			return err
		}
	}

	deletedGroups, err := l.dbEngine.RemoveGroupsFromUser(userID)
	if err != nil {
		return err
	}
	l.logger.Info().Msgf("deleted groups: %v", deletedGroups)

	if len(deletedGroups) > 0 {
		// remove rocket chat
		if err := l.rocketChatEngine.RemoveUseFromGroups(userID, deletedGroups); err != nil {
			return err
		}
	}

	return nil
}
