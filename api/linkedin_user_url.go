package controller

import (
	"context"
	"fmt"
	"strings"

	"github.com/phassans/frolleague/common"
	"github.com/phassans/frolleague/engines"
)

type (
	linkedInUserURLRequest struct {
		LinkedInUserID engines.LinkedInUserID `json:"linkedInUserId"`
		LinkedInURL    engines.LinkedInURL    `json:"linkedInURL"`
	}

	linkedInUserURLResponse struct {
		Ok bool `json:"ok"`
	}

	linkedInUserURLEndpoint struct{}
)

var linkedInUserURL postEndpoint = linkedInUserURLEndpoint{}

func (r linkedInUserURLEndpoint) Execute(ctx context.Context, rtr *router, requestI interface{}) (interface{}, error) {
	request := requestI.(linkedInUserURLRequest)
	if err := r.Validate(requestI); err != nil {
		return loginResponse{}, err
	}

	err := rtr.engines.UpdateUserWithLinkedInURL(request.LinkedInUserID, request.LinkedInURL)
	if err != nil {
		return linkedInUserURLResponse{Ok: false}, err
	}
	return linkedInUserURLResponse{Ok: true}, nil
}

func (r linkedInUserURLEndpoint) Validate(request interface{}) error {
	input := request.(linkedInUserURLRequest)
	if strings.TrimSpace(string(input.LinkedInUserID)) == "" {
		return common.ValidationError{Message: fmt.Sprint("user me failed, missing fields!")}
	}
	return nil
}

func (r linkedInUserURLEndpoint) GetPath() string {
	return "/linkedin/link"
}

func (r linkedInUserURLEndpoint) HTTPRequest() interface{} {
	return linkedInUserURLRequest{}
}

func (r linkedInUserURLEndpoint) GetMessage(err error) string {
	return ""
}
