package engines

import (
	"fmt"
	"strings"

	"github.com/phassans/frolleague/clients/linkedin"

	"github.com/phassans/frolleague/clients/phantom"
	"github.com/phassans/frolleague/clients/rocket"
	"github.com/rs/zerolog"
)

type (
	userEngine struct {
		rClient  rocket.Client
		pClient  phantom.Client
		dbEngine DatabaseEngine
		logger   zerolog.Logger
	}

	UserEngine interface {
		SignUp(Username, Password, LinkedInURL) (User, error)
		Login(Username, Password) (User, error)
		Refresh(UserID) error
		ChangePassword(UserID, Password) error
		DeleteUser(UserID) error
		UpdateUserWithImage(UserID, ImageLink) error
		UpdateUserWithLinkedInURL(id linkedin.UserID, url LinkedInURL) error

		GetUserChatGroups(UserID) ([]GroupWithStatus, error)
		ToggleUserGroup(UserID, Group, bool) error
	}
)

func NewUserEngine(rClient rocket.Client, pClient phantom.Client, dbEngine DatabaseEngine, logger zerolog.Logger) (UserEngine, error) {
	return &userEngine{
		rClient,
		pClient,
		dbEngine,
		logger,
	}, nil
}

func (u *userEngine) SignUp(username Username, password Password, linkedInURL LinkedInURL) (User, error) {
	// add user to db
	var userId UserID
	var err error

	// add user
	userId, err = u.dbEngine.AddUser(username, password, linkedInURL)
	if err != nil {
		return User{}, err
	}

	profile, groups, err := u.getAndProcessUserProfile(linkedInURL, userId)
	if err != nil {
		return User{}, err
	}

	if err := u.BootstrapRocketUser(username, password, profile, groups); err != nil {
		return User{}, err
	}

	return User{UserID: userId}, nil
}

func (u *userEngine) BootstrapRocketUser(username Username, password Password, profile phantom.Profile, groups []Group) error {
	// user
	userID, err := u.createUserIfNotExist(username, password, profile)
	if err != nil {
		return err
	}
	u.logger.Info().Msgf("user: %s added to rocket!", userID)

	// groups
	groupIDs, err := u.createGroupsIfNotExist(groups)
	if err != nil {
		return err
	}
	u.logger.Info().Msgf("groups: %s added to rocket!", groupIDs)

	// addUserToGroup
	for _, groupID := range groupIDs {
		resp, err := u.rClient.AddUserToGroup(rocket.AddUserToGroupRequest{RoomId: groupID, UserId: userID})
		if err != nil {
			u.logger.Info().Msgf("AddUserToGroup failed, user: %s and group: %s", userID, groupID)
			return err
		}
		if !resp.Success {
			u.logger.Info().Msgf("user: %s and group: %s association exists", userID, groupID)
		}
	}

	return nil
}

func (u *userEngine) createUserIfNotExist(username Username, password Password, profile phantom.Profile) (string, error) {
	var userId string
	infoUserResp, err := u.rClient.InfoUser(rocket.InfoUserRequest{Username: string(username)})
	if err != nil {
		if err, ok := err.(rocket.ErrorInvalidUser); !ok {
			return userId, err
		}
	}
	userId = infoUserResp.User.ID

	if infoUserResp.Success == false {
		// create user
		name := fmt.Sprintf("%s %s", profile.User.Firstname, profile.User.LastName)
		email := fmt.Sprintf("%s@gmail.com", username)
		resp, err := u.rClient.CreateUser(rocket.CreateUserRequest{Name: name, Username: string(username), Password: string(password), Email: email})
		if err != nil {
			return userId, err
		}
		userId = resp.User.ID
	}
	return userId, nil
}

func (u *userEngine) createGroupsIfNotExist(groups []Group) ([]string, error) {
	// create groups
	var groupIDs []string
	for _, group := range groups {
		var groupID string
		groupInfo, err := u.rClient.InfoGroup(rocket.InfoGroupRequest{RoomName: string(group)})
		if err != nil {
			if err, ok := err.(rocket.ErrorGroupNotFound); !ok {
				return groupIDs, err
			}
		}

		if groupInfo.Success == false {
			resp, err := u.rClient.CreateGroup(rocket.GroupCreateRequest{Name: string(group)})
			if err != nil {
				if !strings.Contains(err.Error(), "error-duplicate-channel-name") {
					return nil, err
				}
			}
			groupID = resp.Group.ID
		}
		groupID = groupInfo.Group.ID
		groupIDs = append(groupIDs, groupID)
	}

	return groupIDs, nil
}

func (u *userEngine) Refresh(userID UserID) error {
	var userId UserID
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
	}
	return nil
}

func (u *userEngine) ChangePassword(userID UserID, password Password) error {
	var err error

	// getUser
	_, err = u.dbEngine.GetUserByUserID(userID)
	if err != nil {
		return err
	}

	// update user preferences
	if err := u.dbEngine.UpdateUserPassword(userID, password); err != nil {
		return err
	}

	return nil
}

func (u *userEngine) DeleteUser(userID UserID) error {
	var err error

	// getUser
	_, err = u.dbEngine.GetUserByUserID(userID)
	if err != nil {
		return err
	}

	// not deleting user now, will think about it more

	return nil
}

func (u *userEngine) Login(username Username, password Password) (User, error) {
	return u.dbEngine.GetUserByUserNameAndPassword(username, password)
}

func (u *userEngine) GetUserChatGroups(userID UserID) ([]GroupWithStatus, error) {
	return u.dbEngine.GetGroupsWithStatusByUserID(userID)
}

func (u *userEngine) ToggleUserGroup(userID UserID, group Group, status bool) error {
	return u.dbEngine.ToggleUserGroup(userID, group, status)
}

func (u *userEngine) getAndProcessUserProfile(linkedInURL LinkedInURL, userId UserID) (phantom.Profile, []Group, error) {
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
	if err := u.dbEngine.UpdateUserWithNameAndReference(FirstName(profile.User.Firstname), LastName(profile.User.LastName), FileName(profile.FileName), userId); err != nil {
		return phantom.Profile{}, nil, err
	}

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

func (u *userEngine) UpdateUserWithImage(userID UserID, imageLink ImageLink) error {
	return u.dbEngine.UpdateUserWithImage(userID, imageLink)
}

func (u *userEngine) UpdateUserWithLinkedInURL(id linkedin.UserID, url LinkedInURL) error {
	fmt.Printf("linkedInURL: %s", url)
	return nil
}
