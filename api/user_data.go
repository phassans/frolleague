package controller

import (
	"context"
	"fmt"
	"strings"

	"github.com/phassans/frolleague/common"
	"github.com/phassans/frolleague/engines"
)

type (
	fetchUserDataRequest struct {
		UserID  engines.UserID `json:"userId"`
		School  bool           `json:"school,omitempty"`
		Company bool           `json:"company,omitempty"`
	}

	fetchUserDataResponse struct {
		Companies []engines.Company `json:"companies,omitempty"`
		Schools   []engines.School  `json:"schools,omitempty"`
	}

	fetchUserEducationEndpoint struct{}
)

var fetchUser postEndpoint = fetchUserEducationEndpoint{}

func (r fetchUserEducationEndpoint) Execute(ctx context.Context, rtr *router, requestI interface{}) (interface{}, error) {
	request := requestI.(fetchUserDataRequest)
	if err := r.Validate(requestI); err != nil {
		return nil, err
	}

	if request.Company {
		comps, err := rtr.engines.GetCompaniesByUserID(request.UserID)
		return fetchUserDataResponse{Companies: comps}, err
	} else if request.School {
		schools, err := rtr.engines.GetSchoolsByUserID(request.UserID)
		return fetchUserDataResponse{Schools: schools}, err
	}

	result := fetchUserDataResponse{}
	return result, nil
}

func (r fetchUserEducationEndpoint) Validate(request interface{}) error {
	input := request.(fetchUserDataRequest)
	if strings.TrimSpace(string(input.UserID)) == "" {
		return common.ValidationError{Message: fmt.Sprint("fetch user education failed, missing fields")}
	}

	if input.School == false && input.Company == false {
		return common.ValidationError{Message: fmt.Sprint("fetch user education failed, missing fields. Either company or school should set to true")}
	}

	return nil
}

func (r fetchUserEducationEndpoint) GetPath() string {
	return "/user/fetch"
}

func (r fetchUserEducationEndpoint) HTTPRequest() interface{} {
	return fetchUserDataRequest{}
}

func (r fetchUserEducationEndpoint) HTTPResult() interface{} {
	return fetchUserDataResponse{}
}

func (r fetchUserEducationEndpoint) GetMessage(err error) string {
	// just add a success message
	msg := ""
	if err != nil {
		msg = "failed to add user to group!"
	} else {
		msg = "user added to group success!"
	}
	return msg
}
