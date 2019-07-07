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
	linkedInUserCrawlRequest struct {
		LinkedInUserID engines.LinkedInUserID `json:"linkedInUserId"`
	}

	linkedInUserCrawlResponse struct {
		Profile phantom.Profile `json:"profile,omitempty"`
		Ok      bool            `json:"ok"`
	}

	linkedInUserCrawlEndpoint struct{}
)

var linkedInUserCrawl postEndpoint = linkedInUserCrawlEndpoint{}

func (r linkedInUserCrawlEndpoint) Execute(ctx context.Context, rtr *router, requestI interface{}) (interface{}, error) {
	request := requestI.(linkedInUserCrawlRequest)
	if err := r.Validate(requestI); err != nil {
		return linkedInUserCrawlResponse{}, err
	}

	profile, err := rtr.engines.CrawlUserProfile(request.LinkedInUserID)
	if err != nil {
		return linkedInUserCrawlResponse{Ok: false, Profile: profile}, err
	}
	return linkedInUserCrawlResponse{Ok: true, Profile: profile}, nil
}

func (r linkedInUserCrawlEndpoint) Validate(request interface{}) error {
	input := request.(linkedInUserCrawlRequest)
	if strings.TrimSpace(string(input.LinkedInUserID)) == "" {
		return common.ValidationError{Message: fmt.Sprint("user me failed, missing fields!")}
	}
	return nil
}

func (r linkedInUserCrawlEndpoint) GetPath() string {
	return "/linkedin/crawl"
}

func (r linkedInUserCrawlEndpoint) HTTPRequest() interface{} {
	return linkedInUserCrawlRequest{}
}

func (r linkedInUserCrawlEndpoint) GetMessage(err error) string {
	return ""
}
