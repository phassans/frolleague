package controller

import (
	"context"
	"fmt"
	"strings"

	"github.com/phassans/frolleague/common"
)

type (
	linkedInCookieRequest struct {
		UserName string `json:userName`
		Cookie   string `json:"cookie"`
	}

	linkedInCookieResponse struct {
		Ok bool `json:"ok"`
	}

	linkedInCookieEndpoint struct{}
)

var linkedInCookie postEndpoint = linkedInCookieEndpoint{}

func (r linkedInCookieEndpoint) Execute(ctx context.Context, rtr *router, requestI interface{}) (interface{}, error) {
	request := requestI.(linkedInCookieRequest)
	if err := r.Validate(requestI); err != nil {
		return linkedInCookieResponse{}, err
	}

	if err := rtr.engines.UpdateCookie(request.UserName, request.Cookie); err != nil {
		return linkedInCookieResponse{Ok: false}, err
	}

	return linkedInCookieResponse{Ok: true}, nil
}

func (r linkedInCookieEndpoint) Validate(request interface{}) error {
	input := request.(linkedInCookieRequest)
	if strings.TrimSpace(string(input.Cookie)) == "" ||
		strings.TrimSpace(string(input.UserName)) == "" {
		return common.ValidationError{Message: fmt.Sprint("cookie update failed!")}
	}
	return nil
}

func (r linkedInCookieEndpoint) GetPath() string {
	return "/linkedin/cookie"
}

func (r linkedInCookieEndpoint) HTTPRequest() interface{} {
	return linkedInCookieRequest{}
}

func (r linkedInCookieEndpoint) GetMessage(err error) string {
	return ""
}
