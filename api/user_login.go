package controller

import (
	"context"
	"fmt"
	"strings"

	"github.com/phassans/frolleague/common"
	"github.com/phassans/frolleague/engines"
)

type (
	loginRequest struct {
		UserName engines.Username `json:"userName"`
		Password engines.Password `json:"password,omitempty"`
	}

	loginResponse struct {
		User    engines.User `json:"user,omitempty"`
		Error   *APIError    `json:"error,omitempty"`
		Message string       `json:"message,omitempty"`
	}

	loginEndpoint struct{}
)

var login postEndpoint = loginEndpoint{}

func (r loginEndpoint) Execute(ctx context.Context, rtr *router, requestI interface{}) (interface{}, error) {
	request := requestI.(loginRequest)
	if err := r.Validate(requestI); err != nil {
		return loginResponse{}, err
	}

	user, err := rtr.engines.Login(request.UserName, request.Password)
	result := loginResponse{Error: NewAPIError(err), User: user, Message: r.GetMessage(err)}
	return result, err
}

func (r loginEndpoint) Validate(request interface{}) error {
	input := request.(loginRequest)
	if strings.TrimSpace(string(input.UserName)) == "" ||
		strings.TrimSpace(string(input.Password)) == "" {
		return common.ValidationError{Message: fmt.Sprint("login failed, missing fields")}
	}
	return nil
}

func (r loginEndpoint) GetPath() string {
	return "/login"
}

func (r loginEndpoint) HTTPRequest() interface{} {
	return loginRequest{}
}

func (r loginEndpoint) GetMessage(err error) string {
	// just add a success message
	msg := ""
	if err != nil {
		msg = "user login failed!"
	} else {
		msg = "login success!"
	}
	return msg
}
