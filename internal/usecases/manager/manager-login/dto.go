package managerlogin

import (
	"github.com/evgeniy-krivenko/chat-service/internal/types"
	"github.com/evgeniy-krivenko/chat-service/internal/validator"
)

type Request struct {
	Login    string `validate:"min=3,max=32"`
	Password string `validate:"min=3,max=32"`
}

func (r Request) Validate() error {
	return validator.Validator.Struct(r)
}

type Response struct {
	Token     string
	ClientID  types.UserID
	FirstName string
	LastName  string
}
