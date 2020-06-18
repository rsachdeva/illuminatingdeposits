package web

import (
	"fmt"

	"github.com/pkg/errors"
)

// FieldError is used to indicate an error with a specific request field.
type FieldError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

// ErrorResponse is the form used for API responses from failures in the API.
type ErrorResponse struct {
	Error  string       `json:"error"`
	Fields []FieldError `json:"fields,omitempty"`
}

// ErrorRequest is used to pass an error during the request through the
// application with web specific context.
type ErrorRequest struct { // Think of this as web.ErrorField - it is just a custome error uses cause error message for ErrorField() implementation
	Err    error
	Status int
	Fields []FieldError
}

// NewRequestError wraps a provided error with an HTTP status code. This
// function should be used when handlers encounter expected errors.
func NewRequestError(err error, status int) error {
	return &ErrorRequest{err, status, nil}
}

// ErrorField implements the error interface. It uses the default message of the
// wrapped error. This is what will be shown in the services' logs.
func (err *ErrorRequest) Error() string {
	fmt.Println("\terr.Err.Error() is", err.Err.Error())
	return err.Err.Error()
}

// shutdown is a type used to help with the graceful termination of the service.
type shutdown struct {
	Message string
}

// Error is the implementation of the error interface.
func (s *shutdown) Error() string {
	return s.Message
}

// NewShutdownError returns an error that causes the framework to signal
// a graceful shutdown.
func NewShutdownError(message string) error {
	return &shutdown{message}
}

// IsShutdown checks to see if the shutdown error is contained
// in the specified error value.
func IsShutdown(err error) bool {
	if _, ok := errors.Cause(err).(*shutdown); ok {
		return true
	}
	return false
}
