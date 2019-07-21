package engines

import (
	"github.com/phassans/frolleague/clients/phantom"
	"github.com/phassans/frolleague/common"
	"github.com/rs/zerolog"
)

type (
	userEngine struct {
		pClient  phantom.Client
		dbEngine DatabaseEngine
		logger   zerolog.Logger
	}

	UserEngine interface {
		// old
		Refresh(UserID) error

		GetUserChatGroups(UserID) (AllGroups, error)
		ToggleUserGroup(UserID, Group, bool) error

		// new
		GetSchoolsByUserID(userID UserID) ([]School, error)
		GetCompaniesByUserID(userID UserID) ([]Company, error)
	}
)

func NewUserEngine(pClient phantom.Client, dbEngine DatabaseEngine, logger zerolog.Logger) (UserEngine, error) {
	return &userEngine{
		pClient,
		dbEngine,
		logger,
	}, nil
}

func (u *userEngine) Refresh(userID UserID) error {
	/*var userId UserID
	var err error

	// getUser
	user, err := u.dbEngine.GetUserByUserID(userID)
	if err != nil {
		return err
	}

	profile, groups, err := u.getAndProcessUserProfile(user.LinkedInURL, userId)
	if err != nil {
		return err
	}

	if err := u.BootstrapRocketUser(user.Username, "", profile, groups); err != nil {
		return err
	}*/
	return nil
}

func (u *userEngine) GetUserChatGroups(userID UserID) (AllGroups, error) {

	var allGroups AllGroups

	groups, err := u.dbEngine.GetGroupsWithStatusByUserID(userID)
	if err != nil {
		return allGroups, err
	}

	var companyGroups []GroupWithStatus
	var schoolGroups []GroupWithStatus
	for _, group := range groups {
		if group.GroupSource == "company" {
			companyGroups = append(companyGroups, group)
		}
		if group.GroupSource == "school" {
			schoolGroups = append(schoolGroups, group)
		}
	}
	allGroups.CompanyGroups = companyGroups
	allGroups.SchoolGroups = schoolGroups

	return allGroups, nil
}

func (u *userEngine) ToggleUserGroup(userID UserID, group Group, status bool) error {
	_, err := u.dbEngine.GetUserByID(LinkedInUserID(userID))
	if err != nil {
		return err
	}

	groups, err := u.dbEngine.GetGroupsByUserID(userID)
	if err != nil {
		return err
	}

	isValidGroup := false
	for _, g := range groups {
		if g.Group == group {
			isValidGroup = true
		}
	}
	if !isValidGroup {
		return common.ErrorNotExist{"user group doesnt exist!"}
	}

	return u.dbEngine.ToggleUserGroup(userID, group, status)
}

func (u *userEngine) getAndProcessUserProfile(linkedInURL LinkedInURL, userId UserID) (phantom.Profile, []GroupInfo, error) {
	// get userProfile
	profile, err := u.pClient.GetUserProfile(string(linkedInURL), false)
	if err != nil {
		return phantom.Profile{}, nil, err
	}

	if err := u.addUserToSchools(profile, userId); err != nil {
		return phantom.Profile{}, nil, err
	}

	if err := u.addUserToCompanies(profile, userId); err != nil {
		return phantom.Profile{}, nil, err
	}

	// update user preferences
	/*if err := u.dbEngine.UpdateUserWithNameAndReference(FirstName(profile.User.Firstname), LastName(profile.User.LastName), FileName(profile.FileName), userId); err != nil {
		return phantom.Profile{}, nil, err
	}*/

	// update user preferences
	groups, err := u.dbEngine.AddGroupsToUser(userId)
	if err != nil {
		return phantom.Profile{}, nil, err
	}

	return profile, groups, nil
}

func (u *userEngine) addUserToSchools(profile phantom.Profile, userID UserID) error {
	for _, school := range profile.Schools {
		schoolID, err := u.dbEngine.AddSchoolIfNotPresent(SchoolName(school.SchoolName), Degree(school.Degree), FieldOfStudy(school.FieldOfStudy))
		if err != nil {
			return err
		}

		if err := u.dbEngine.AddUserToSchool(userID, schoolID, FromYear(school.FromYear), ToYear(school.ToYear)); err != nil {
			return err
		}
	}
	return nil
}

func (u *userEngine) addUserToCompanies(profile phantom.Profile, userID UserID) error {
	for _, company := range profile.Companies {
		companyID, err := u.dbEngine.AddCompanyIfNotPresent(CompanyName(company.CompanyName), Location(company.Location))
		if err != nil {
			return err
		}

		if err := u.dbEngine.AddUserToCompany(userID, companyID, Title(company.Title), FromYear(company.FromYear), ToYear(company.ToYear)); err != nil {
			return err
		}
	}

	return nil
}

func (u *userEngine) GetSchoolsByUserID(userID UserID) ([]School, error) {
	return u.dbEngine.GetSchoolsByUserID(userID)
}

func (u *userEngine) GetCompaniesByUserID(userID UserID) ([]Company, error) {
	return u.dbEngine.GetCompaniesByUserID(userID)
}
