package getuserprofile

import (
	"github.com/evgeniy-krivenko/chat-service/internal/types"
	"github.com/evgeniy-krivenko/chat-service/internal/validator"
)

type Request struct {
	ID     types.RequestID `validate:"required"`
	UserID types.UserID    `validate:"required"`
}

func (r Request) Validate() error {
	return validator.Validator.Struct(r)
}

type Response struct {
	UserID    types.UserID
	FirstName string
	LastName  string
}
