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
	userSchoolsUpdateRequest struct {
		LinkedInId engines.UserID   `json:"linkedInUserId"`
		Schools    []phantom.School `json:"schools"`
	}

	userSchoolsUpdateResponse struct {
		Ok bool `json:"ok"`
	}

	userSchoolsUpdateEndpoint struct{}
)

var userSchoolsUpdate postEndpoint = userSchoolsUpdateEndpoint{}

func (r userSchoolsUpdateEndpoint) Execute(ctx context.Context, rtr *router, requestI interface{}) (interface{}, error) {
	request := requestI.(userSchoolsUpdateRequest)
	if err := r.Validate(requestI); err != nil {
		return userSchoolsUpdateResponse{}, err
	}

	err := rtr.engines.UpdateUserSchools(request.LinkedInId, request.Schools)
	if err != nil {
		return userSchoolsUpdateResponse{Ok: false}, err
	}
	return userSchoolsUpdateResponse{Ok: true}, nil
}

func (r userSchoolsUpdateEndpoint) Validate(request interface{}) error {
	input := request.(userSchoolsUpdateRequest)
	if strings.TrimSpace(string(input.LinkedInId)) == "" {
		return common.ValidationError{Message: fmt.Sprint("user linkedIn URL Update failed!")}
	}
	return nil
}

func (r userSchoolsUpdateEndpoint) GetPath() string {
	return "/user/schools/update"
}

func (r userSchoolsUpdateEndpoint) HTTPRequest() interface{} {
	return userSchoolsUpdateRequest{}
}

func (r userSchoolsUpdateEndpoint) GetMessage(err error) string {
	return ""
}
