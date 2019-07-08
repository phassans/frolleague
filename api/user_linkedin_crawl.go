package controller

import (
	"context"
	"fmt"
	"strings"

	"github.com/phassans/frolleague/clients/phantom"

	"github.com/phassans/frolleague/common"
	"github.com/phassans/frolleague/engines"
)

type (
	userLinkedInCrawlRequest struct {
		UserID engines.LinkedInUserID `json:"userId"`
	}

	userLinkedInCrawlResponse struct {
		Profile phantom.Profile `json:"profile,omitempty"`
		Ok      bool            `json:"ok"`
	}

	userLinkedInCrawlEndpoint struct{}
)

var userLinkedInCrawl postEndpoint = userLinkedInCrawlEndpoint{}

func (r userLinkedInCrawlEndpoint) Execute(ctx context.Context, rtr *router, requestI interface{}) (interface{}, error) {
	request := requestI.(userLinkedInCrawlRequest)
	if err := r.Validate(requestI); err != nil {
		return userLinkedInCrawlResponse{}, err
	}

	profile, err := rtr.engines.CrawlUserProfile(request.UserID)
	if err != nil {
		return userLinkedInCrawlResponse{Ok: false, Profile: profile}, err
	}
	return userLinkedInCrawlResponse{Ok: true, Profile: profile}, nil
}

func (r userLinkedInCrawlEndpoint) Validate(request interface{}) error {
	input := request.(userLinkedInCrawlRequest)
	if strings.TrimSpace(string(input.UserID)) == "" {
		return common.ValidationError{Message: fmt.Sprint("user me failed, missing fields!")}
	}
	return nil
}

func (r userLinkedInCrawlEndpoint) GetPath() string {
	return "/user/linkedin/crawl"
}

func (r userLinkedInCrawlEndpoint) HTTPRequest() interface{} {
	return userLinkedInCrawlRequest{}
}

func (r userLinkedInCrawlEndpoint) GetMessage(err error) string {
	return ""
}
