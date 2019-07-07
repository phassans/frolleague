package controller

import (
	"context"
	"fmt"
	"strings"

	"github.com/phassans/frolleague/clients/linkedin"

	"github.com/phassans/frolleague/common"
	"github.com/phassans/frolleague/engines"
)

type (
	linkedInLogInRequest struct {
		linkedin.AuthCode    `json:"authCode,omitempty"`
		linkedin.AccessToken `json:"accessToken,omitempty"`
	}

	linkedInLogInResponse struct {
		LinkedInId engines.LinkedInId `json:"linkedInId"`
		Message    string             `json:"message,omitempty"`
	}

	linkedInLogInEndpoint struct{}
)

var linkedInLogIn postEndpoint = linkedInLogInEndpoint{}

func (r linkedInLogInEndpoint) Execute(ctx context.Context, rtr *router, requestI interface{}) (interface{}, error) {
	request := requestI.(linkedInLogInRequest)
	if err := r.Validate(requestI); err != nil {
		return loginResponse{}, err
	}

	resp, err := rtr.engines.LogIn(request.AuthCode, request.AccessToken)
	result := linkedInLogInResponse{LinkedInId: engines.LinkedInId(resp.ID)}
	return result, err
}

func (r linkedInLogInEndpoint) Validate(request interface{}) error {
	input := request.(linkedInLogInRequest)
	if strings.TrimSpace(string(input.AuthCode)) == "" && strings.TrimSpace(string(input.AccessToken)) == "" {
		return common.ValidationError{Message: fmt.Sprint("user log in failed, missing fields!")}
	}
	return nil
}

func (r linkedInLogInEndpoint) GetPath() string {
	return "/linkedin/login"
}

func (r linkedInLogInEndpoint) HTTPRequest() interface{} {
	return linkedInLogInRequest{}
}

func (r linkedInLogInEndpoint) GetMessage(err error) string {
	return ""
}
