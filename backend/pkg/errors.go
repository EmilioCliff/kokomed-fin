package pkg

import (
	"errors"
	"fmt"
	"net/http"
)

const (
	ALREADY_EXISTS_ERROR  = "already_exists"
	INTERNAL_ERROR        = "internal"
	INVALID_ERROR         = "invalid"
	NOT_FOUND_ERROR       = "not_found"
	NOT_IMPLEMENTED_ERROR = "not_implemented"
	AUTHENTICATION_ERROR  = "authentication"
)

type Error struct {
	Code    string
	Message string
}

func Errorf(code string, format string, args ...any) *Error {
	return &Error{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}
}

func ErrorCode(err error) string {
	var e *Error

	if err == nil {
		return ""
	} else if errors.As(err, &e) {
		return e.Code
	}

	return INTERNAL_ERROR
}

func ErrorMessage(err error) string {
	var e *Error

	if err == nil {
		return ""
	} else if errors.As(err, &e) {
		return e.Message
	}

	return "Internal error."
}

func ErrorToStatusCode(err error) int {
	switch ErrorCode(err) {
	case ALREADY_EXISTS_ERROR:
		return http.StatusConflict
	case INTERNAL_ERROR:
		return http.StatusInternalServerError
	case INVALID_ERROR:
		return http.StatusBadRequest
	case NOT_FOUND_ERROR:
		return http.StatusNotFound
	case NOT_IMPLEMENTED_ERROR:
		return http.StatusNotImplemented
	case AUTHENTICATION_ERROR:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}

// ALREADY_EXISTS_ERROR  = "already_exists"
// INTERNAL_ERROR        = "internal"
// INVALID_ERROR         = "invalid"
// NOT_FOUND_ERROR       = "not_found"
// NOT_IMPLEMENTED_ERROR = "not_implemented"
// AUTHENTICATION_ERROR  = "authentication"

// Error implements the error interface. Not used by the application otherwise.
func (e *Error) Error() string {
	return fmt.Sprintf("error: code=%s message=%s", e.Code, e.Message)
}
