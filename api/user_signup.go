package controller

import (
	"context"
	"fmt"
	"strings"

	"github.com/phassans/frolleague/common"
	"github.com/phassans/frolleague/engines"
)

type (
	signUpRequest struct {
		UserName    engines.Username    `json:"userName"`
		Password    engines.Password    `json:"password,omitempty"`
		LinkedInURL engines.LinkedInURL `json:"linkedInURL"`
	}

	signUpResponse struct {
		UserId  engines.UserID `json:"userId"`
		Error   *APIError      `json:"error,omitempty"`
		Message string         `json:"message,omitempty"`
	}

	signUpEndpoint struct{}
)

var signUp postEndpoint = signUpEndpoint{}

func (r signUpEndpoint) Execute(ctx context.Context, rtr *router, requestI interface{}) (interface{}, error) {
	request := requestI.(signUpRequest)

	if err := r.Validate(requestI); err != nil {
		return nil, err
	}

	user, err := rtr.engines.SignUp(request.UserName, request.Password, request.LinkedInURL)
	result := signUpResponse{Error: NewAPIError(err), UserId: user.UserID, Message: r.GetMessage(err)}
	return result, err
}

func (r signUpEndpoint) Validate(request interface{}) error {
	input := request.(signUpRequest)
	if strings.TrimSpace(string(input.UserName)) == "" ||
		strings.TrimSpace(string(input.Password)) == "" {
		return common.ValidationError{Message: fmt.Sprint("signUp failed, missing fields")}
	}
	return nil
}

func (r signUpEndpoint) GetPath() string {
	return "/signup"
}

func (r signUpEndpoint) HTTPRequest() interface{} {
	return signUpRequest{}
}

func (r signUpEndpoint) GetMessage(err error) string {
	// just add a success message
	msg := ""
	if err != nil {
		msg = "failed signing up user!"
	} else {
		msg = "user signup success!"
	}
	return msg
}
