package controller

import (
	"context"
	"fmt"
	"net/url"

	"github.com/phassans/banana/helper"
	"github.com/phassans/frolleague/engines"
)

type (
	userLinkedInURLGETEndpoint struct{}
)

var userLinkedInURLGET getEndPoint = userLinkedInURLGETEndpoint{}

func (r userLinkedInURLGETEndpoint) Do(ctx context.Context, rtr *router, values url.Values) (interface{}, error) {
	if values.Get("userId") == "" {
		return nil, helper.ValidationError{Message: fmt.Sprint("missing userId!")}
	}
	return rtr.engines.GetLinkedInURL(engines.LinkedInUserID(values.Get("userId")))
}

func (r userLinkedInURLGETEndpoint) GetPath() string {
	return "/user/linkedin/url"
}
