package engines

import (
	"fmt"
	"strings"

	"github.com/phassans/frolleague/clients/rocket"
	"github.com/rs/zerolog"
)

type (
	rocketChatEngine struct {
		rClient  rocket.Client
		dbEngine DatabaseEngine
		logger   zerolog.Logger
	}

	RocketChatEngine interface {
		SetUpRocketChatForUser(user UserID) error
		CreateGroupsIfNotExist(groups []GroupInfo) ([]string, error)
		AddUserToGroups(userID UserID, groups []GroupInfo) error
		RemoveUseFromGroups(userID UserID, groups []GroupInfo) error
		GetUserIDFromName(username Username) (string, bool, error)
	}
)

func NewRocketChatEngine(rClient rocket.Client, dbEngine DatabaseEngine, logger zerolog.Logger) RocketChatEngine {
	return &rocketChatEngine{
		rClient,
		dbEngine,
		logger,
	}
}

func (r *rocketChatEngine) SetUpRocketChatForUser(userID UserID) error {
	user, err := r.dbEngine.GetUserByID(LinkedInUserID(userID))
	if err != nil {
		return err
	}
	userName := Username(fmt.Sprintf("%s.%s", user.FirstName, user.LastName))

	//get Groups
	groups, err := r.dbEngine.GetGroupsByUserIDAndStatus(userID, true)
	if err != nil {
		return err
	}

	err = r.BootstrapRocketUser(userName, Password("password"), groups)
	if err != nil {
		return err
	}

	return nil
}

func (r *rocketChatEngine) BootstrapRocketUser(username Username, password Password, groups []GroupInfo) error {
	// user
	userID, err := r.createUserIfNotExist(username, password)
	if err != nil {
		return err
	}
	r.logger.Info().Msgf("user: %s added to rocket!", userID)

	// groups
	groupIDs, err := r.CreateGroupsIfNotExist(groups)
	if err != nil {
		return err
	}
	r.logger.Info().Msgf("groups: %s added to rocket!", groupIDs)

	// addUserToGroups
	if err := r.doAddUserToRCGroups(userID, groupIDs); err != nil {
		return err
	}

	return nil
}

func (r *rocketChatEngine) GetUserIDFromName(username Username) (string, bool, error) {
	var userId string
	infoUserResp, err := r.rClient.InfoUser(rocket.InfoUserRequest{Username: string(username)})
	if err != nil {
		if err, ok := err.(rocket.ErrorInvalidUser); !ok {
			return userId, infoUserResp.Success, err
		}
	}
	userId = infoUserResp.User.ID
	return userId, infoUserResp.Success, nil
}

func (r *rocketChatEngine) createUserIfNotExist(username Username, password Password) (string, error) {
	userId, success, err := r.GetUserIDFromName(username)
	if err != nil {
		return userId, err
	}

	if success == false {
		// create user
		email := fmt.Sprintf("%s@gmail.com", username)
		resp, err := r.rClient.CreateUser(rocket.CreateUserRequest{
			Name:     string(username),
			Username: string(username),
			Password: string(password),
			Email:    email,
		})
		if err != nil {
			return userId, err
		}
		userId = resp.User.ID
	}
	return userId, nil
}

func (r *rocketChatEngine) doAddUserToRCGroups(userID string, groupIDs []string) error {
	for _, groupID := range groupIDs {
		resp, err := r.rClient.AddUserToGroup(rocket.AddUserToGroupRequest{RoomId: groupID, UserId: userID})
		if err != nil {
			r.logger.Info().Msgf("AddUserToGroup failed, user: %s and group: %s", userID, groupID)
			return err
		}
		if !resp.Success {
			r.logger.Info().Msgf("user: %s and group: %s association exists", userID, groupID)
		}
	}

	return nil
}

func (r *rocketChatEngine) doRemoveUserFromRCGroups(userID string, groupIDs []string) error {
	for _, groupID := range groupIDs {
		resp, err := r.rClient.RemoveUserFromGroup(rocket.RemoveUserFromGroupRequest{RoomId: groupID, UserId: userID})
		if err != nil {
			r.logger.Info().Msgf("AddUserToGroup failed, user: %s and group: %s", userID, groupID)
			return err
		}
		if !resp.Success {
			r.logger.Info().Msgf("user: %s and group: %s association exists", userID, groupID)
		}
	}

	return nil
}

func (r *rocketChatEngine) AddUserToGroups(userID UserID, groups []GroupInfo) error {
	return r.doUserGroup(userID, groups, "add")
}

func (r *rocketChatEngine) RemoveUseFromGroups(userID UserID, groups []GroupInfo) error {
	return r.doUserGroup(userID, groups, "delete")
}

func (r *rocketChatEngine) doUserGroup(userID UserID, groups []GroupInfo, operation string) error {
	// getUserInfo
	user, err := r.dbEngine.GetUserByID(LinkedInUserID(userID))
	if err != nil {
		return err
	}
	userName := Username(fmt.Sprintf("%s.%s", user.FirstName, user.LastName))

	// getRCUserID
	rcUserID, status, err := r.GetUserIDFromName(userName)
	if err != nil {
		return err
	}

	// user exists
	if status == true {
		// createGroups
		groupIDs, err := r.CreateGroupsIfNotExist(groups)
		if err != nil {
			return err
		}

		if operation == "add" {
			// AddUsersToGroups
			if err := r.doAddUserToRCGroups(rcUserID, groupIDs); err != nil {
				return err
			}
		} else if operation == "delete" {
			// RemoveUsersToGroups
			if err := r.doRemoveUserFromRCGroups(rcUserID, groupIDs); err != nil {
				return err
			}
		}
	}

	return nil
}

func (u *rocketChatEngine) CreateGroupsIfNotExist(groups []GroupInfo) ([]string, error) {
	// create groups
	var groupIDs []string
	for _, group := range groups {
		var groupID string
		groupInfo, err := u.rClient.InfoGroup(rocket.InfoGroupRequest{RoomName: string(group.Group)})
		if err != nil {
			if err, ok := err.(rocket.ErrorGroupNotFound); !ok {
				return groupIDs, err
			}
		}
		groupID = groupInfo.Group.ID

		if groupInfo.Success == false {
			resp, err := u.rClient.CreateGroup(rocket.GroupCreateRequest{Name: string(group.Group)})
			if err != nil {
				if !strings.Contains(err.Error(), "error-duplicate-channel-name") {
					return nil, err
				}
			}
			groupID = resp.Group.ID
		}
		groupIDs = append(groupIDs, groupID)
	}

	return groupIDs, nil
}
