package controller

import (
	"context"
	"fmt"
	"strings"

	"github.com/phassans/frolleague/common"
	"github.com/phassans/frolleague/engines"
)

type (
	userLinkedInURLRequest struct {
		UserID      engines.LinkedInUserID `json:"userId"`
		LinkedInURL engines.LinkedInURL    `json:"linkedInURL"`
	}

	userLinkedInURLResponse struct {
		Ok bool `json:"ok"`
	}

	userLinkedInURLEndpoint struct{}
)

var userLinkedInURL postEndpoint = userLinkedInURLEndpoint{}

func (r userLinkedInURLEndpoint) Execute(ctx context.Context, rtr *router, requestI interface{}) (interface{}, error) {
	request := requestI.(userLinkedInURLRequest)
	if err := r.Validate(requestI); err != nil {
		return userLinkedInURLResponse{}, err
	}

	err := rtr.engines.UpdateUserWithLinkedInURL(request.UserID, request.LinkedInURL)
	if err != nil {
		return userLinkedInURLResponse{Ok: false}, err
	}
	return userLinkedInURLResponse{Ok: true}, nil
}

func (r userLinkedInURLEndpoint) Validate(request interface{}) error {
	input := request.(userLinkedInURLRequest)
	if strings.TrimSpace(string(input.UserID)) == "" ||
		strings.TrimSpace(string(input.LinkedInURL)) == "" {
		return common.ValidationError{Message: fmt.Sprint("user linkedIn URL Update failed!")}
	}
	return nil
}

func (r userLinkedInURLEndpoint) GetPath() string {
	return "/user/linkedin/url"
}

func (r userLinkedInURLEndpoint) HTTPRequest() interface{} {
	return userLinkedInURLRequest{}
}

func (r userLinkedInURLEndpoint) GetMessage(err error) string {
	return ""
}
