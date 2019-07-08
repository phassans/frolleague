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
	userCompaniesUpdateRequest struct {
		UserID    engines.UserID    `json:"userId"`
		Companies []phantom.Company `json:"companies"`
	}

	userCompaniesUpdateResponse struct {
		Ok bool `json:"ok"`
	}

	userCompaniesUpdateEndpoint struct{}
)

var userCompaniesUpdate postEndpoint = userCompaniesUpdateEndpoint{}

func (r userCompaniesUpdateEndpoint) Execute(ctx context.Context, rtr *router, requestI interface{}) (interface{}, error) {
	request := requestI.(userCompaniesUpdateRequest)
	if err := r.Validate(requestI); err != nil {
		return userCompaniesUpdateResponse{}, err
	}

	err := rtr.engines.UpdateUserCompanies(request.UserID, request.Companies)
	if err != nil {
		return userCompaniesUpdateResponse{Ok: false}, err
	}
	return userCompaniesUpdateResponse{Ok: true}, nil
}

func (r userCompaniesUpdateEndpoint) Validate(request interface{}) error {
	input := request.(userCompaniesUpdateRequest)
	if strings.TrimSpace(string(input.UserID)) == "" {
		return common.ValidationError{Message: fmt.Sprint("user linkedIn URL Update failed!")}
	}
	return nil
}

func (r userCompaniesUpdateEndpoint) GetPath() string {
	return "/user/companies/update"
}

func (r userCompaniesUpdateEndpoint) HTTPRequest() interface{} {
	return userCompaniesUpdateRequest{}
}

func (r userCompaniesUpdateEndpoint) GetMessage(err error) string {
	return ""
}
