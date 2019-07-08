package controller

import (
	"context"
	"fmt"
	"strings"

	"github.com/phassans/frolleague/clients/linkedin"

	"github.com/phassans/frolleague/common"
)

type (
	userLinkedInMeRequest struct {
		UserID linkedin.UserID `json:"userId"`
	}

	userLinkedInMeResponse struct {
		Ok        bool   `json:"isUserLogin"`
		FirstName string `json:"firstName,omitempty"`
		LastName  string `json:"lastName,omitempty"`
		Message   string `json:"message,omitempty"`
	}

	userLinkedInMeEndpoint struct{}
)

var userLinkedInMe postEndpoint = userLinkedInMeEndpoint{}

func (r userLinkedInMeEndpoint) Execute(ctx context.Context, rtr *router, requestI interface{}) (interface{}, error) {
	request := requestI.(userLinkedInMeRequest)
	if err := r.Validate(requestI); err != nil {
		return userLinkedInMeResponse{}, err
	}

	meresp, ok, err := rtr.engines.GetMe(request.UserID)
	if !ok {
		return userLinkedInMeResponse{Ok: ok}, err
	}
	return userLinkedInMeResponse{Ok: ok, FirstName: meresp.FirstName.Localized.EnUS, LastName: meresp.LastName.Localized.EnUS}, err
}

func (r userLinkedInMeEndpoint) Validate(request interface{}) error {
	input := request.(userLinkedInMeRequest)
	if strings.TrimSpace(string(input.UserID)) == "" {
		return common.ValidationError{Message: fmt.Sprint("user linkedIn me failed, missing fields!")}
	}
	return nil
}

func (r userLinkedInMeEndpoint) GetPath() string {
	return "/user/linkedin/me"
}

func (r userLinkedInMeEndpoint) HTTPRequest() interface{} {
	return userLinkedInMeRequest{}
}

func (r userLinkedInMeEndpoint) GetMessage(err error) string {
	return ""
}
