package controller

import (
	"context"
	"fmt"
	"strings"

	"github.com/phassans/frolleague/common"
	"github.com/phassans/frolleague/engines"
)

type (
	userLinkedInURLPOSTRequest struct {
		UserID      engines.LinkedInUserID `json:"userId"`
		LinkedInURL engines.LinkedInURL    `json:"linkedInURL"`
	}

	userLinkedInURLPOSTResponse struct {
		Ok bool `json:"ok"`
	}

	userLinkedInURLPOSTEndpoint struct{}
)

var userLinkedInURLPOST postEndpoint = userLinkedInURLPOSTEndpoint{}

func (r userLinkedInURLPOSTEndpoint) Execute(ctx context.Context, rtr *router, requestI interface{}) (interface{}, error) {
	request := requestI.(userLinkedInURLPOSTRequest)
	if err := r.Validate(requestI); err != nil {
		return userLinkedInURLPOSTResponse{}, err
	}

	err := rtr.engines.UpdateUserWithLinkedInURL(request.UserID, request.LinkedInURL)
	if err != nil {
		return userLinkedInURLPOSTResponse{Ok: false}, err
	}
	return userLinkedInURLPOSTResponse{Ok: true}, nil
}

func (r userLinkedInURLPOSTEndpoint) Validate(request interface{}) error {
	input := request.(userLinkedInURLPOSTRequest)
	if strings.TrimSpace(string(input.UserID)) == "" ||
		strings.TrimSpace(string(input.LinkedInURL)) == "" {
		return common.ValidationError{Message: fmt.Sprint("user linkedIn URL Update failed!")}
	}
	return nil
}

func (r userLinkedInURLPOSTEndpoint) GetPath() string {
	return "/user/linkedin/url"
}

func (r userLinkedInURLPOSTEndpoint) HTTPRequest() interface{} {
	return userLinkedInURLPOSTRequest{}
}

func (r userLinkedInURLPOSTEndpoint) GetMessage(err error) string {
	return ""
}
