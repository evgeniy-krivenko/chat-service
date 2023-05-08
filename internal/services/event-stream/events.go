package eventstream

import (
	"fmt"
	"github.com/evgeniy-krivenko/chat-service/internal/types"
	"github.com/evgeniy-krivenko/chat-service/internal/validator"
	"time"
)

type Event interface {
	eventMarker()
	Validate() error
}

type event struct{}         //
func (*event) eventMarker() {}

// NewMessageEvent is a signal about the appearance of a new message in the chat.
type NewMessageEvent struct {
	event

	EventID     types.EventID   `validate:"required"`
	RequestID   types.RequestID `validate:"required"`
	ChatID      types.ChatID    `validate:"required"`
	MessageID   types.MessageID `validate:"required"`
	AuthorID    types.UserID    `validate:"omitempty"`
	CreatedAt   time.Time       `validate:"omitempty"`
	MessageBody string          `validate:"omitempty"`
	IsService   bool
}

func (e NewMessageEvent) Validate() error {
	if err := validator.Validator.Struct(e); err != nil {
		return fmt.Errorf("validate NewMessageEvent: %v", err)
	}

	return nil
}

func NewNewMessageEvent(
	eventID types.EventID,
	reqID types.RequestID,
	chatID types.ChatID,
	msgID types.MessageID,
	authorID types.UserID,
	createdAt time.Time,
	msgBody string,
	isService bool,
) *NewMessageEvent {
	return &NewMessageEvent{
		event:       event{},
		MessageBody: msgBody,
		EventID:     eventID,
		RequestID:   reqID,
		ChatID:      chatID,
		MessageID:   msgID,
		AuthorID:    authorID,
		CreatedAt:   createdAt,
		IsService:   isService,
	}
}

type MessageSentEvent struct {
	event

	EventID   types.EventID   `validate:"required"`
	RequestID types.RequestID `validate:"required"`
	MessageID types.MessageID `validate:"required"`
}

func (m MessageSentEvent) Validate() error {
	if err := validator.Validator.Struct(m); err != nil {
		return fmt.Errorf("validate MessageSentEvent: %v", err)
	}

	return nil
}

func NewMessageSentEvent(
	eventID types.EventID,
	reqID types.RequestID,
	msgID types.MessageID,
) *MessageSentEvent {
	return &MessageSentEvent{
		EventID:   eventID,
		RequestID: reqID,
		MessageID: msgID,
	}
}
