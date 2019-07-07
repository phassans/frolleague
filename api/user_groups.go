package controller

import (
	"context"
	"fmt"

	"github.com/phassans/frolleague/common"
	"github.com/phassans/frolleague/engines"
)

type (
	userGroupsRequest struct {
		UserID engines.UserID `json:"userId"`
	}

	userGroupsResponse struct {
		Groups  []engines.GroupWithStatus `json:"groups,omitempty"`
		Error   *APIError                 `json:"error,omitempty"`
		Message string                    `json:"message,omitempty"`
	}

	userGroupsEndpoint struct{}
)

var userGroups postEndpoint = userGroupsEndpoint{}

func (r userGroupsEndpoint) Execute(ctx context.Context, rtr *router, requestI interface{}) (interface{}, error) {
	request := requestI.(userGroupsRequest)
	if err := r.Validate(requestI); err != nil {
		return loginResponse{}, err
	}

	groups, err := rtr.engines.GetUserChatGroups(request.UserID)
	result := userGroupsResponse{Error: NewAPIError(err), Groups: groups}
	return result, err
}

func (r userGroupsEndpoint) Validate(request interface{}) error {
	input := request.(userGroupsRequest)
	if input.UserID == 0 {
		return common.ValidationError{Message: fmt.Sprint("invalid UserID for userGroups")}
	}
	return nil
}

func (r userGroupsEndpoint) GetPath() string {
	return "/usergroups"
}

func (r userGroupsEndpoint) HTTPRequest() interface{} {
	return userGroupsRequest{}
}

func (r userGroupsEndpoint) GetMessage(err error) string {
	// just add a success message
	msg := ""
	if err != nil {
		msg = "failed to add user to group!"
	} else {
		msg = "user added to group success!"
	}
	return msg
}
