package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/phassans/frolleague/clients/linkedin"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/phassans/frolleague/common"
)

// APIError is a HTTP result error.
type APIError struct {
	Code         int    `json:"code"`
	ErrorMessage string `json:"error"`
}

// NewAPIError returns a new result error.
func NewAPIError(err error) *APIError {
	if err == nil {
		return nil
	}

	return &APIError{Code: GetErrorStatus(err), ErrorMessage: err.Error()}
}

func (e *APIError) Error() string {
	return fmt.Sprintf("api error [%d]: %s", e.Code, e.Error())
}

// Send writes result error into response.json.
func (e *APIError) Send(w http.ResponseWriter) error {
	w.WriteHeader(e.Code)
	return json.NewEncoder(w).Encode(e)
}

// GetErrorStatus returns the HTTP status code for a given error type.
func GetErrorStatus(err error) int {
	if e, ok := err.(*APIError); err == nil || ok && e == nil {
		return http.StatusOK
	}
	switch err := err.(type) {
	case *APIError:
		return err.Code
	case common.UserError:
		return http.StatusBadRequest
	case common.ErrorUserNotExist:
		return http.StatusBadRequest
	case common.DuplicateSignUp:
		return http.StatusBadRequest
	case common.ValidationError:
		return http.StatusBadRequest
	case linkedin.APIPostError:
		return err.Code
	case linkedin.APIGetError:
		return err.Status
	default:
		return http.StatusInternalServerError
	}
}

// IsHardError designates whether an error is hard enough to open hystrix circuit breaker
func IsHardError(err error) bool {
	switch err.(type) {
	case hystrix.CircuitError:
		return true
	default:
		return false
	}
}
