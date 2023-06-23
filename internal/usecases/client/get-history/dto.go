package gethistory

import (
	"errors"
	"time"

	messagesrepo "github.com/evgeniy-krivenko/chat-service/internal/repositories/messages"
	"github.com/evgeniy-krivenko/chat-service/internal/types"
	"github.com/evgeniy-krivenko/chat-service/internal/validator"
)

type Request struct {
	ID       types.RequestID `validate:"required"`
	ClientID types.UserID    `validate:"required"`
	PageSize int             `validate:"omitempty,gte=10,lte=100"`
	Cursor   string          `validate:"omitempty,base64url"`
}

func (r Request) Validate() error {
	if r.Cursor == "" && r.PageSize == 0 {
		return errors.New("neither cursor nor pagesize specified")
	}
	if r.Cursor != "" && r.PageSize != 0 {
		return errors.New("cursor and pagesize specified")
	}

	return validator.Validator.Struct(r)
}

type Response struct {
	Messages   []Message
	NextCursor string
}

type Message struct {
	ID                  types.MessageID
	AuthorID            types.UserID
	AuthorName          string
	Body                string
	CreatedAt           time.Time
	IsBlocked           bool
	IsVisibleForManager bool
	IsReceived          bool
	IsService           bool
}

func adaptMessage(m messagesrepo.Message) Message {
	return Message{
		ID:                  m.ID,
		AuthorID:            m.AuthorID,
		AuthorName:          m.AuthorFirstName,
		Body:                m.Body,
		CreatedAt:           m.CreatedAt,
		IsBlocked:           m.IsBlocked,
		IsVisibleForManager: false,
		IsReceived:          m.IsVisibleForManager && !m.IsBlocked,
		IsService:           m.IsService,
	}
}
