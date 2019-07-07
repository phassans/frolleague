package common

import (
	"encoding/json"
	"fmt"
)

// ValidationError ...
type ValidationError struct {
	Message string `json:"message,omitempty"`
}

func (v ValidationError) Error() string {
	b, _ := json.Marshal(v)
	return fmt.Sprintf("validation error: %s", string(b))
}

// UserError ...
type UserError struct {
	Message string `json:"message,omitempty"`
}

func (v UserError) Error() string {
	b, _ := json.Marshal(v)
	return fmt.Sprintf("%s", string(b))
}

// ErrorUserNotExist ...
type ErrorUserNotExist struct {
	Message string `json:"message,omitempty"`
}

func (v ErrorUserNotExist) Error() string {
	b, _ := json.Marshal(v)
	return fmt.Sprintf("%s", string(b))
}

// DuplicateSignUp ...
type DuplicateSignUp struct {
	Username    string `json:"username,omitempty"`
	LinkedInURL string `json:"linkedInURL,omitempty"`
	Message     string `json:"message,omitempty"`
}

func (e DuplicateSignUp) Error() string {
	b, _ := json.Marshal(e)
	return fmt.Sprintf("%s", string(b))
}

// DatabaseError ...
type DatabaseError struct {
	DBError string `json:"dbError,omitempty"`
}

func (e DatabaseError) Error() string {
	b, _ := json.Marshal(e)
	return fmt.Sprintf("database error: %s", string(b))
}

// ValidationError ...
type LocationError struct {
	Message string `json:"message,omitempty"`
}

func (l LocationError) Error() string {
	b, _ := json.Marshal(l)
	return fmt.Sprintf("location error: %s", string(b))
}

// LinkedInError ...
type LinkedInError struct {
	Message string `json:"message,omitempty"`
}

func (l LinkedInError) Error() string {
	b, _ := json.Marshal(l)
	return fmt.Sprintf("%s", string(b))
}
