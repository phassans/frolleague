package controller

import (
	"context"
	"fmt"
	"strings"

	"github.com/phassans/frolleague/common"
	"github.com/phassans/frolleague/engines"
)

type (
	userChatSetUpRequest struct {
		UserID engines.UserID `json:"userId"`
	}

	userChatSetUpResponse struct {
		Ok bool `json:"ok"`
	}

	userChatSetUpEndpoint struct{}
)

var userChatSetUp postEndpoint = userChatSetUpEndpoint{}

func (r userChatSetUpEndpoint) Execute(ctx context.Context, rtr *router, requestI interface{}) (interface{}, error) {
	request := requestI.(userChatSetUpRequest)
	if err := r.Validate(requestI); err != nil {
		return userChatSetUpResponse{}, err
	}

	err := rtr.engines.SetUpRocketChatForUser(request.UserID)
	if err != nil {
		return userChatSetUpResponse{Ok: false}, err
	}
	return userChatSetUpResponse{Ok: true}, nil
}

func (r userChatSetUpEndpoint) Validate(request interface{}) error {
	input := request.(userChatSetUpRequest)
	if strings.TrimSpace(string(input.UserID)) == "" {
		return common.ValidationError{Message: fmt.Sprint("userID empty for setting up user!")}
	}
	return nil
}

func (r userChatSetUpEndpoint) GetPath() string {
	return "/user/chat/setup"
}

func (r userChatSetUpEndpoint) HTTPRequest() interface{} {
	return userChatSetUpRequest{}
}

func (r userChatSetUpEndpoint) GetMessage(err error) string {
	return ""
}
