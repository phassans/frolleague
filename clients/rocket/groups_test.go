package rocket

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClient_CreateGroup(t *testing.T) {
	newRocketChatClient(t)
	{
		// create Group with no name
		createGroupResp, err := rClient.CreateGroup(GroupCreateRequest{})
		require.Error(t, err)
		require.False(t, createGroupResp.Success)

		// create Group with invalid name
		createGroupResp, err = rClient.CreateGroup(GroupCreateRequest{testGroupName + " 1234"})
		require.NoError(t, err)
		require.False(t, createGroupResp.Success)

		// create Group success
		createGroupResp, err = rClient.CreateGroup(GroupCreateRequest{testGroupName})
		require.NoError(t, err)
		require.True(t, createGroupResp.Success)
		groupID := createGroupResp.Group.ID

		// create Group duplicate
		createGroupResp, err = rClient.CreateGroup(GroupCreateRequest{testGroupName})
		require.Error(t, err)
		require.False(t, createGroupResp.Success)

		// delete Group
		deleteGroupResp, err := rClient.DeleteGroup(DeleteGroupRequest{groupID})
		require.NoError(t, err)
		require.True(t, deleteGroupResp.Success)
	}
}

func TestClient_InfoGroup(t *testing.T) {
	newRocketChatClient(t)
	{
		// get Group info invalid group
		infoGroupResp, err := rClient.InfoGroup(InfoGroupRequest{testGroupName + "invalid"})
		require.Error(t, err)
		require.False(t, infoGroupResp.Success)

		// get Group info empty group
		infoGroupResp, err = rClient.InfoGroup(InfoGroupRequest{})
		require.Error(t, err)
		require.False(t, infoGroupResp.Success)

		// create Group success
		createGroupResp, err := rClient.CreateGroup(GroupCreateRequest{testGroupName})
		require.NoError(t, err)
		require.True(t, createGroupResp.Success)

		// get Group info
		infoGroupResp, err = rClient.InfoGroup(InfoGroupRequest{testGroupName})
		require.NoError(t, err)
		require.True(t, infoGroupResp.Success)

		// delete Group
		deleteGroupResp, err := rClient.DeleteGroup(DeleteGroupRequest{createGroupResp.Group.ID})
		require.NoError(t, err)
		require.True(t, deleteGroupResp.Success)
	}
}

func TestClient_DeleteGroup(t *testing.T) {
	newRocketChatClient(t)
	{
		// delete Group invalid
		deleteGroupResp, err := rClient.DeleteGroup(DeleteGroupRequest{testGroupName + "invalid"})
		require.Error(t, err)
		require.False(t, deleteGroupResp.Success)

		// delete Group empty
		deleteGroupResp, err = rClient.DeleteGroup(DeleteGroupRequest{})
		require.Error(t, err)
		require.False(t, deleteGroupResp.Success)

		// create Group success
		createGroupResp, err := rClient.CreateGroup(GroupCreateRequest{testGroupName})
		require.NoError(t, err)
		require.True(t, createGroupResp.Success)

		// delete Group
		deleteGroupResp, err = rClient.DeleteGroup(DeleteGroupRequest{createGroupResp.Group.ID})
		require.NoError(t, err)
		require.True(t, deleteGroupResp.Success)
	}
}

func TestClient_SetTypeGroup(t *testing.T) {
	newRocketChatClient(t)
	{
	}
}
