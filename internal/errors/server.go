package errors

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

const internalErrorMsg = "something went wrong"

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
	if errSrv := new(ServerError); errors.As(err, &errSrv) {
		return errSrv.Code, errSrv.Message, errSrv.Error()
	}

	if errEcho := new(echo.HTTPError); errors.As(err, &errEcho) {
		return errEcho.Code, errEcho.Message.(string), errEcho.Error()
	}

	return http.StatusInternalServerError, internalErrorMsg, err.Error()
}
