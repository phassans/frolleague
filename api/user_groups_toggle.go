package controller

import (
	"context"
	"fmt"

	"github.com/phassans/frolleague/common"
	"github.com/phassans/frolleague/engines"
)

type (
	userGroupToggleRequest struct {
		UserID engines.UserID `json:"userId"`
		Group  engines.Group  `json:"group"`
		Status bool           `json:"status"`
	}

	userGroupToggleResponse struct {
		Request userGroupToggleRequest `json:"request,omitempty"`
		Error   *APIError              `json:"error,omitempty"`
		Message string                 `json:"message,omitempty"`
	}

	userGroupToggleEndpoint struct{}
)

var userGroupToggle postEndpoint = userGroupToggleEndpoint{}

func (r userGroupToggleEndpoint) Execute(ctx context.Context, rtr *router, requestI interface{}) (interface{}, error) {
	request := requestI.(userGroupToggleRequest)
	if err := r.Validate(requestI); err != nil {
		return loginResponse{}, err
	}

	err := rtr.engines.ToggleUserGroup(request.UserID, request.Group, request.Status)
	result := userGroupToggleResponse{Request: request, Error: NewAPIError(err), Message: r.GetMessage(err)}
	return result, err
}

func (r userGroupToggleEndpoint) Validate(request interface{}) error {
	input := request.(userGroupToggleRequest)
	if input.UserID == 0 || input.Group == "" {
		return common.ValidationError{Message: fmt.Sprint("invalid userId or group to toggle usergroup")}
	}
	return nil
}

func (r userGroupToggleEndpoint) GetPath() string {
	return "/usergroup/toggle"
}

func (r userGroupToggleEndpoint) HTTPRequest() interface{} {
	return userGroupToggleRequest{}
}

func (r userGroupToggleEndpoint) GetMessage(err error) string {
	// just add a success message
	msg := ""
	if err != nil {
		msg = "failed to toggle user group!"
	} else {
		msg = "user group toggle success!"
	}
	return msg
}
