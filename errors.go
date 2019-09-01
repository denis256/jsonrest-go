package jsonrest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Error creates an error that will be rendered directly to the client.
func Error(status int, code, message string) error {
	return &httpError{
		Status:  status,
		Code:    code,
		Message: message,
	}
}

// BadRequest returns an HTTP 400 Bad Request error with a custom error message.
func BadRequest(msg string) error {
	return Error(http.StatusBadRequest, "bad_request", msg)
}

// UnauthorizedError returns an HTTP 401 Unauthorized error with a custom error
// message.
func UnauthorizedError(msg string) error {
	return Error(http.StatusUnauthorized, "unauthorized", msg)
}

// NotFound returns an HTTP 404 Not Found error with a custom error message.
func NotFound(msg string) error {
	return Error(http.StatusNotFound, "not_found", msg)
}

// UnprocessableEntity returns an HTTP 422 UnprocessableEntity error with a
// custom error message.
func UnprocessableEntity(msg string) error {
	return Error(http.StatusUnprocessableEntity, "unprocessable_entity", msg)
}

// unknownError is returned for an internal server error.
var unknownError = &httpError{
	Code:    "unknown_error",
	Message: "an unknown error occurred",
	Status:  500,
}

// httpError is an error that will be rendered to the client.
type httpError struct {
	Code    string
	Message string
	Details []string
	Status  int
}

// MarshalJSON implements the json.Marshaler interface.
func (err *httpError) MarshalJSON() ([]byte, error) {
	var wp struct {
		Error struct {
			Code    string   `json:"code"`
			Message string   `json:"message"`
			Details []string `json:"details,omitempty"`
		} `json:"error"`
	}
	wp.Error.Code = err.Code
	wp.Error.Message = err.Message
	wp.Error.Details = err.Details
	return json.Marshal(wp)
}

// Error implements the error interface.
func (err *httpError) Error() string {
	return fmt.Sprintf("jsonrest: %v: %v", err.Code, err.Message)
}

// translateError coerces err into an httpError that can be marshaled directly
// to the client.
func translateError(err error, dumpInternalError bool) *httpError {
	httpErr, ok := err.(*httpError)
	if !ok {
		httpErr = &(*unknownError) // shallow copy
		if dumpInternalError {
			httpErr.Details = dumpError(err)
		}
	}
	return httpErr
}

// dumpError formats the error suitable for viewing in a JSON response for local
// debugging.
func dumpError(err error) []string {
	s := fmt.Sprintf("%+v", err)           // stringify
	s = strings.Replace(s, "\t", "  ", -1) // tabs to spaces
	return strings.Split(s, "\n")          // split on newline
}