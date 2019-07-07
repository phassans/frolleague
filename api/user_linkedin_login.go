package controller

import (
	"context"
	"fmt"
	"strings"

	"github.com/phassans/frolleague/clients/linkedin"

	"github.com/phassans/frolleague/common"
)

type (
	userLinkedInLogInRequest struct {
		linkedin.AuthCode    `json:"authCode,omitempty"`
		linkedin.AccessToken `json:"accessToken,omitempty"`
	}

	userLinkedInLogInResponse struct {
		LinkedInId string `json:"linkedInId"`
		Message    string `json:"message,omitempty"`
	}

	userLinkedInLogInEndpoint struct{}
)

var userLinkedInLogIn postEndpoint = userLinkedInLogInEndpoint{}

func (r userLinkedInLogInEndpoint) Execute(ctx context.Context, rtr *router, requestI interface{}) (interface{}, error) {
	request := requestI.(userLinkedInLogInRequest)
	if err := r.Validate(requestI); err != nil {
		return userLinkedInLogInResponse{}, err
	}

	resp, err := rtr.engines.LogIn(request.AuthCode, request.AccessToken)
	result := userLinkedInLogInResponse{LinkedInId: resp.ID}
	return result, err
}

func (r userLinkedInLogInEndpoint) Validate(request interface{}) error {
	input := request.(userLinkedInLogInRequest)
	if strings.TrimSpace(string(input.AuthCode)) == "" && strings.TrimSpace(string(input.AccessToken)) == "" {
		return common.ValidationError{Message: fmt.Sprint("user log in failed, missing fields!")}
	}
	return nil
}

func (r userLinkedInLogInEndpoint) GetPath() string {
	return "/user/linkedin/login"
}

func (r userLinkedInLogInEndpoint) HTTPRequest() interface{} {
	return userLinkedInLogInRequest{}
}

func (r userLinkedInLogInEndpoint) GetMessage(err error) string {
	return ""
}
