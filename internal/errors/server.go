package errors

import (
	"errors"
	"fmt"

	"github.com/labstack/echo/v4"
)

// ServerError is used to return custom error codes to client.
type ServerError struct {
	Code    int
	Message string
	cause   error
}

func NewServerError[T ~int](code T, msg string, err error) *ServerError {
	return &ServerError{Message: msg, Code: int(code), cause: err}
}

func (s *ServerError) Error() string {
	return fmt.Sprintf("%s: %v", s.Message, s.cause)
}

func (s *ServerError) Unwrap() error {
	return s.cause
}

func GetServerErrorCode(err error) int {
	code, _, _ := ProcessServerError(err)
	return code
}

// ProcessServerError tries to retrieve from given error it's code, message and some details.
// For example, that fields can be used to build error response for client.
func ProcessServerError(err error) (code int, msg string, details string) {
	var serverErr *ServerError
	if errors.As(err, &serverErr) {
		code = serverErr.Code
		msg = serverErr.Message
		details = serverErr.Error()
		return
	}

	var echoErr *echo.HTTPError
	if errors.As(err, &echoErr) {
		code = echoErr.Code
		msg = echoErr.Message.(string)
		details = echoErr.Error()
		return
	}

	code = 500
	msg = "something went wrong"
	details = err.Error()
	return
}
