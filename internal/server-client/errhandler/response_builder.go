package errhandler

import (
	clientv1 "github.com/evgeniy-krivenko/chat-service/internal/server-client/v1"
)

type Response struct {
	Error clientv1.Error `json:"error"`
}

var ResponseBuilder = func(code int, msg string, details string) any {
	var d *string
	if details != "" {
		d = &details
	}
	return Response{
		Error: clientv1.Error{
			Code:    code,
			Message: msg,
			Details: d,
		},
	}
}
