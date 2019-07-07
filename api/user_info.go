package controller

import (
	"context"
	"fmt"
	"strings"

	"github.com/phassans/frolleague/common"
	"github.com/phassans/frolleague/engines"
)

type (
	userInfoRequest struct {
		UserID  engines.UserID `json:"userId"`
		School  bool           `json:"school,omitempty"`
		Company bool           `json:"company,omitempty"`
	}

	userInfoResponse struct {
		Companies []engines.Company `json:"companies,omitempty"`
		Schools   []engines.School  `json:"schools,omitempty"`
	}

	userInfoEndPoint struct{}
)

var userInfo postEndpoint = userInfoEndPoint{}

func (r userInfoEndPoint) Execute(ctx context.Context, rtr *router, requestI interface{}) (interface{}, error) {
	request := requestI.(userInfoRequest)
	if err := r.Validate(requestI); err != nil {
		return nil, err
	}

	if request.Company {
		comps, err := rtr.engines.GetCompaniesByUserID(request.UserID)
		return userInfoResponse{Companies: comps}, err
	} else if request.School {
		schools, err := rtr.engines.GetSchoolsByUserID(request.UserID)
		return userInfoResponse{Schools: schools}, err
	}

	result := userInfoResponse{}
	return result, nil
}

func (r userInfoEndPoint) Validate(request interface{}) error {
	input := request.(userInfoRequest)
	if strings.TrimSpace(string(input.UserID)) == "" {
		return common.ValidationError{Message: fmt.Sprint("fetch user education failed, missing fields")}
	}

	if input.School == false && input.Company == false {
		return common.ValidationError{Message: fmt.Sprint("fetch user education failed, missing fields. Either company or school should set to true")}
	}

	return nil
}

func (r userInfoEndPoint) GetPath() string {
	return "/user/info"
}

func (r userInfoEndPoint) HTTPRequest() interface{} {
	return userInfoRequest{}
}

func (r userInfoEndPoint) HTTPResult() interface{} {
	return userInfoResponse{}
}

func (r userInfoEndPoint) GetMessage(err error) string {
	// just add a success message
	msg := ""
	if err != nil {
		msg = "failed to add user to group!"
	} else {
		msg = "user added to group success!"
	}
	return msg
}
