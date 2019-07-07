package rocket

import "github.com/rs/zerolog"

type (
	client struct {
		baseURL string
		logger  zerolog.Logger
	}

	Client interface {
		// baseMethod for making a post call
		DoPost(interface{}, string, AdminCredentials) ([]byte, error)

		// base methods
		InitClient(username string, password string) error
		GetAdminCredentials() AdminCredentials

		// user login, mostly to obtain accessToken & userId
		Login(UserLoginRequest) (UserLoginResponse, error)
		InfoUser(InfoUserRequest) (InfoUserResponse, error)
		CreateUser(CreateUserRequest) (CreateUserResponse, error)
		DeleteUser(DeleteUserRequest) (DeleteUserResponse, error)

		// creates a new Group
		CreateGroup(GroupCreateRequest) (GroupCreateResponse, error)
		DeleteGroup(DeleteGroupRequest) (DeleteGroupResponse, error)
		InfoGroup(InfoGroupRequest) (InfoGroupResponse, error)
		SetTypeGroup(SetTypeGroupRequest) (SetTypeGroupResponse, error)

		AddUserToGroup(AddUserToGroupRequest) (AddUserToGroupResponse, error)
		RemoveUserFromGroup(RemoveUserFromGroupRequest) (RemoveUserFromGroupResponse, error)
	}
)

// NewRocketClient returns a new rocket client
func NewRocketClient(baseURL string, logger zerolog.Logger) Client {
	return &client{baseURL, logger}
}
