package rocket

import (
	"fmt"
	"net/http"
)

const (
	ErrorUserParamNotProvidedType = "error-user-param-not-provided"
	ErrorInvalidUserType          = "error-invalid-user"

	ErrorGroupNotFoundType         = "error-room-not-found"
	ErrorGroupParamNotProvidedType = "error-room-param-not-provided"
	ErrorInvalidGroupNameType      = "error-invalid-room-name"
	ErrorDuplicateGroupNameType    = "error-duplicate-channel-name"
)

type ErrHTTP struct {
	Request  *http.Request
	Response *http.Response
}

func (e ErrHTTP) Error() string {
	return fmt.Sprintf("%s %s returned %s", e.Request.Method, e.Request.URL, e.Response.Status)
}

type ErrorDuplicateUserName struct {
	Username string
}

func (e ErrorDuplicateUserName) Error() string {
	return fmt.Sprintf("username %s exists", e.Username)
}

type ErrorUserCreation struct {
	ErrorMsg string
}

func (e ErrorUserCreation) Error() string {
	return fmt.Sprintf("UserCreationError %s", e.ErrorMsg)
}

type ErrorRequiredParam struct {
	ErrorMsg string
}

func (e ErrorRequiredParam) Error() string {
	return fmt.Sprintf("RequiredParamError %s", e.ErrorMsg)
}

type ErrorInvalidUser struct {
	ErrorMsg string
}

func (e ErrorInvalidUser) Error() string {
	return fmt.Sprintf("InvalidUser %s", e.ErrorMsg)
}

type ErrorDuplicateGroupName struct {
	GroupName string
}

func (e ErrorDuplicateGroupName) Error() string {
	return fmt.Sprintf("group %s exists", e.GroupName)
}

type ErrorGroupNotFound struct {
	GroupName string
}

func (e ErrorGroupNotFound) Error() string {
	return fmt.Sprintf("group %s not found", e.GroupName)
}

type ErrorInvalidGroupName struct {
	GroupName string
}

func (e ErrorInvalidGroupName) Error() string {
	return fmt.Sprintf("group name %s invalid", e.GroupName)
}
