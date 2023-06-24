package getmanagerprofile

import (
	"github.com/evgeniy-krivenko/chat-service/internal/types"
	"github.com/evgeniy-krivenko/chat-service/internal/validator"
)

type Request struct {
	ID        types.RequestID `validate:"required"`
	ManagerID types.UserID    `validate:"required"`
}

func (r Request) Validate() error {
	return validator.Validator.Struct(r)
}

type Response struct {
	ManagerID types.UserID
	FirstName string
	LastName  string
}
