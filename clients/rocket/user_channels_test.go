package rocket

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClient_AddUserToGroup(t *testing.T) {
	newRocketChatClient(t)
	{
		// add usertoGroup invalid group
		addUserToGroupResp, err := rClient.AddUserToGroup(AddUserToGroupRequest{"1234", "foobar123"})
		require.Error(t, err)
		require.False(t, addUserToGroupResp.Success)

		// add user to Group empty group
		addUserToGroupResp, err = rClient.AddUserToGroup(AddUserToGroupRequest{"", "foobar123"})
		require.Error(t, err)
		require.False(t, addUserToGroupResp.Success)

		// create Group success
		createGroupResp, err := rClient.CreateGroup(GroupCreateRequest{testGroupName})
		require.NoError(t, err)
		require.True(t, createGroupResp.Success)

		// add usertoGroup invalid userID
		addUserToGroupResp, err = rClient.AddUserToGroup(AddUserToGroupRequest{createGroupResp.Group.ID, "foobar123"})
		require.Error(t, err)
		require.False(t, addUserToGroupResp.Success)

		// add usertoGroup empty userID
		addUserToGroupResp, err = rClient.AddUserToGroup(AddUserToGroupRequest{createGroupResp.Group.ID, ""})
		require.Error(t, err)
		require.False(t, addUserToGroupResp.Success)

		// create user success
		resp, err := rClient.CreateUser(CreateUserRequest{testName, testUserEmail, testUsername, testUserPassword})
		require.NoError(t, err)
		require.True(t, resp.Success)

		// add usertoGroup success
		addUserToGroupResp, err = rClient.AddUserToGroup(AddUserToGroupRequest{createGroupResp.Group.ID, resp.User.ID})
		require.NoError(t, err)
		require.True(t, addUserToGroupResp.Success)

		// remove user from Group
		removeUserFromGroupResp, err := rClient.RemoveUserFromGroup(RemoveUserFromGroupRequest{createGroupResp.Group.ID, resp.User.ID})
		require.NoError(t, err)
		require.Equal(t, true, removeUserFromGroupResp.Success)

		// delete user success
		deleteUserResp, err := rClient.DeleteUser(DeleteUserRequest{UserId: resp.User.ID})
		require.NoError(t, err)
		require.True(t, deleteUserResp.Success)

		// delete Group
		deleteGroupResp, err := rClient.DeleteGroup(DeleteGroupRequest{createGroupResp.Group.ID})
		require.NoError(t, err)
		require.True(t, deleteGroupResp.Success)
	}
}

func TestClient_RemoveUserFromGroup(t *testing.T) {
	newRocketChatClient(t)
	{
		// remove usertoGroup invalid group
		removeUserFromGroupResp, err := rClient.RemoveUserFromGroup(RemoveUserFromGroupRequest{"1234", "foobar123"})
		require.Error(t, err)
		require.False(t, removeUserFromGroupResp.Success)

		// remove user to Group empty group
		removeUserFromGroupResp, err = rClient.RemoveUserFromGroup(RemoveUserFromGroupRequest{"", "foobar123"})
		require.Error(t, err)
		require.False(t, removeUserFromGroupResp.Success)

		// create Group success
		createGroupResp, err := rClient.CreateGroup(GroupCreateRequest{testGroupName})
		require.NoError(t, err)
		require.True(t, createGroupResp.Success)

		// create user success
		resp, err := rClient.CreateUser(CreateUserRequest{testName, testUserEmail, testUsername, testUserPassword})
		require.NoError(t, err)
		require.True(t, resp.Success)

		// add usertoGroup success
		addUserToGroupResp, err := rClient.AddUserToGroup(AddUserToGroupRequest{createGroupResp.Group.ID, resp.User.ID})
		require.NoError(t, err)
		require.True(t, addUserToGroupResp.Success)

		// add usertoGroup invalid userID
		removeUserFromGroupResp, err = rClient.RemoveUserFromGroup(RemoveUserFromGroupRequest{createGroupResp.Group.ID, "foobar123"})
		require.Error(t, err)
		require.False(t, removeUserFromGroupResp.Success)

		// add usertoGroup empty userID
		removeUserFromGroupResp, err = rClient.RemoveUserFromGroup(RemoveUserFromGroupRequest{createGroupResp.Group.ID, ""})
		require.Error(t, err)
		require.False(t, removeUserFromGroupResp.Success)

		// remove user from Group
		removeUserFromGroupResp, err = rClient.RemoveUserFromGroup(RemoveUserFromGroupRequest{createGroupResp.Group.ID, resp.User.ID})
		require.NoError(t, err)
		require.Equal(t, true, removeUserFromGroupResp.Success)

		// delete user success
		deleteUserResp, err := rClient.DeleteUser(DeleteUserRequest{UserId: resp.User.ID})
		require.NoError(t, err)
		require.True(t, deleteUserResp.Success)

		// delete Group
		deleteGroupResp, err := rClient.DeleteGroup(DeleteGroupRequest{createGroupResp.Group.ID})
		require.NoError(t, err)
		require.True(t, deleteGroupResp.Success)
	}
}
