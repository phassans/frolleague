package rocket

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClient_CreateUser(t *testing.T) {
	newRocketChatClient(t)
	{
		// invalid email
		resp, err := rClient.CreateUser(CreateUserRequest{testName, "", testUsername + "123", testUserPassword})
		require.Error(t, err)
		require.False(t, resp.Success)

		// invalid username
		resp, err = rClient.CreateUser(CreateUserRequest{testName, testUserEmail, "", testUserPassword})
		require.Error(t, err)
		require.False(t, resp.Success)

		// success
		resp, err = rClient.CreateUser(CreateUserRequest{testName, testUserEmail, testUsername, testUserPassword})
		require.NoError(t, err)
		require.True(t, resp.Success)
		userID := resp.User.ID

		// duplicate user error
		resp, err = rClient.CreateUser(CreateUserRequest{testName, testUserEmail, testUsername, testUserPassword})
		require.Error(t, err)
		require.False(t, resp.Success)

		// delete user error
		deleteUserResp, err := rClient.DeleteUser(DeleteUserRequest{resp.User.ID})
		require.Error(t, err)
		require.False(t, deleteUserResp.Success)

		// delete user error
		deleteUserResp, err = rClient.DeleteUser(DeleteUserRequest{"1234"})
		require.Error(t, err)
		require.False(t, deleteUserResp.Success)

		// delete user success
		deleteUserResp, err = rClient.DeleteUser(DeleteUserRequest{UserId: userID})
		require.NoError(t, err)
		require.True(t, deleteUserResp.Success)
	}
}

func TestClient_InfoUser(t *testing.T) {
	newRocketChatClient(t)
	{
		// success
		resp, err := rClient.CreateUser(CreateUserRequest{testName, testUserEmail, testUsername, testUserPassword})
		require.NoError(t, err)
		require.True(t, resp.Success)
		userID := resp.User.ID

		// get user info error
		infoUserResp, err := rClient.InfoUser(InfoUserRequest{testUsername + "invalid"})
		require.Error(t, err)
		require.False(t, infoUserResp.Success)

		// get user info error
		infoUserResp, err = rClient.InfoUser(InfoUserRequest{})
		require.Error(t, err)
		require.False(t, infoUserResp.Success)

		// get user info success
		infoUserResp, err = rClient.InfoUser(InfoUserRequest{testUsername})
		require.NoError(t, err)
		require.True(t, infoUserResp.Success)

		// delete user success
		deleteUserResp, err := rClient.DeleteUser(DeleteUserRequest{UserId: userID})
		require.NoError(t, err)
		require.True(t, deleteUserResp.Success)
	}
}
