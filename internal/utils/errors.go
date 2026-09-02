package utils

import (
	"fmt"
)

const (
	errorPrefix     = "ERR_"
	ErrUnauthorized = "UNAUTHORIZED"
	ErrForbidden    = "FORBIDDEN"
	ErrNotFound     = "NOT_FOUND"
	ErrBadRequest   = "BAD_REQUEST"
	ErrInternal     = "INTERNAL_SERVER_ERROR"
	ErrService      = "SERVICE_UNAVAILABLE"
	ErrGateway      = "GATEWAY_TIMEOUT"
	ErrTimeout      = "REQUEST_TIMEOUT"
	ErrConflict     = "CONFLICT"
)

var ErrNil = fmt.Errorf("cannot be nil")

type CustomError struct {
	Message string
	Code    string
	Details string
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func NewCustomError(message, code, details string) *CustomError {
	return &CustomError{
		Message: message,
		Code:    code,
		Details: details,
	}
}

func ErrorResponse(code string, err error, details ...string) error {
	var detailStr string
	if len(details) > 0 {
		detailStr = details[0]
	}

	code = errorPrefix + code
	return NewCustomError(err.Error(), code, detailStr)
}
