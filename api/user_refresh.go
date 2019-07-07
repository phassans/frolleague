package controller

import (
	"context"
	"fmt"
	"strings"

	"github.com/phassans/frolleague/common"
	"github.com/phassans/frolleague/engines"
)

type (
	refreshRequest struct {
		UserID engines.UserID `json:"userId"`
	}

	refreshResponse struct {
		refreshRequest
		Error   *APIError `json:"error,omitempty"`
		Message string    `json:"message,omitempty"`
	}

	refreshEndpoint struct{}
)

var refresh postEndpoint = refreshEndpoint{}

func (r refreshEndpoint) Execute(ctx context.Context, rtr *router, requestI interface{}) (interface{}, error) {
	request := requestI.(refreshRequest)

	if err := r.Validate(requestI); err != nil {
		return nil, err
	}

	err := rtr.engines.Refresh(request.UserID)
	result := refreshResponse{refreshRequest: request, Error: NewAPIError(err), Message: r.GetMessage(err)}
	return result, err
}

func (r refreshEndpoint) Validate(request interface{}) error {
	input := request.(refreshRequest)
	if strings.TrimSpace(string(input.UserID)) == "" {
		return common.ValidationError{Message: fmt.Sprint("refresh failed, missing fields")}
	}
	return nil
}

func (r refreshEndpoint) GetPath() string {
	return "/refresh"
}

func (r refreshEndpoint) HTTPRequest() interface{} {
	return refreshRequest{}
}

func (r refreshEndpoint) GetMessage(err error) string {
	// just add a success message
	msg := ""
	if err != nil {
		msg = "failed refreshing user!"
	} else {
		msg = "user refresh success!"
	}
	return msg
}
