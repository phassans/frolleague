package controller

import (
	"context"
	"fmt"
	"strings"

	"github.com/phassans/frolleague/clients/linkedin"

	"github.com/phassans/frolleague/common"
)

type (
	linkedInUserMeRequest struct {
		LinkedInId linkedin.UserID `json:"linkedInUserId"`
	}

	linkedInUserMeResponse struct {
		Ok        bool   `json:"isUserLogin"`
		FirstName string `json:"firstName,omitempty"`
		LastName  string `json:"lastName,omitempty"`
		Message   string `json:"message,omitempty"`
	}

	linkedInUserMeEndpoint struct{}
)

var linkedInUserAuthCode postEndpoint = linkedInUserMeEndpoint{}

func (r linkedInUserMeEndpoint) Execute(ctx context.Context, rtr *router, requestI interface{}) (interface{}, error) {
	request := requestI.(linkedInUserMeRequest)
	if err := r.Validate(requestI); err != nil {
		return loginResponse{}, err
	}

	meresp, ok, err := rtr.engines.GetMe(request.LinkedInId)
	if !ok {
		return linkedInUserMeResponse{Ok: ok}, err
	}
	return linkedInUserMeResponse{Ok: ok, FirstName: meresp.FirstName.Localized.EnUS, LastName: meresp.LastName.Localized.EnUS}, err
}

func (r linkedInUserMeEndpoint) Validate(request interface{}) error {
	input := request.(linkedInUserMeRequest)
	if strings.TrimSpace(string(input.LinkedInId)) == "" {
		return common.ValidationError{Message: fmt.Sprint("user me failed, missing fields!")}
	}
	return nil
}

func (r linkedInUserMeEndpoint) GetPath() string {
	return "/linkedin/me"
}

func (r linkedInUserMeEndpoint) HTTPRequest() interface{} {
	return linkedInUserMeRequest{}
}

func (r linkedInUserMeEndpoint) GetMessage(err error) string {
	return ""
}
