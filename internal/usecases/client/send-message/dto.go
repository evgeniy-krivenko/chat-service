package sendmessage

import (
	"time"

	"github.com/evgeniy-krivenko/chat-service/internal/types"
	"github.com/evgeniy-krivenko/chat-service/internal/validator"
)

type Request struct {
	ID          types.RequestID `validate:"required"`
	ClientID    types.UserID    `validate:"required"`
	MessageBody string          `validate:"min=1,max=2000"`
}

func (r Request) Validate() error {
	return validator.Validator.Struct(r)
}

type Response struct {
	MessageID types.MessageID
	AuthorID  types.UserID
	CreatedAt time.Time
}
