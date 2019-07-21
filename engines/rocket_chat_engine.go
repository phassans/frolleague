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
	groupIDs, err := r.createGroupsIfNotExist(groups)
	if err != nil {
		return err
	}
	r.logger.Info().Msgf("groups: %s added to rocket!", groupIDs)

	// addUserToGroup
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

func (u *rocketChatEngine) createUserIfNotExist(username Username, password Password) (string, error) {
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
		email := fmt.Sprintf("%s@gmail.com", username)
		resp, err := u.rClient.CreateUser(rocket.CreateUserRequest{
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

func (u *rocketChatEngine) createGroupsIfNotExist(groups []GroupInfo) ([]string, error) {
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
