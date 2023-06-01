package eventstream

import (
	"time"

	"github.com/evgeniy-krivenko/chat-service/internal/types"
	"github.com/evgeniy-krivenko/chat-service/internal/validator"
)

//go:generate gonstructor --output=events.gen.go --type=NewMessageEvent --type=MessageSentEvent --type=MessageBlockedEvent --type=NewChatEvent --type=ChatClosedEvent

type Event interface {
	eventMarker()
	Validate() error
}

type event struct{}         //
func (*event) eventMarker() {}

// NewMessageEvent is a signal about the appearance of a new message in the chat.
type NewMessageEvent struct {
	event `gonstructor:"-"`

	EventID     types.EventID   `validate:"required"`
	RequestID   types.RequestID `validate:"required"`
	ChatID      types.ChatID    `validate:"required"`
	MessageID   types.MessageID `validate:"required"`
	AuthorID    types.UserID    `validate:"omitempty"`
	CreatedAt   time.Time       `validate:"required"`
	MessageBody string          `validate:"required,max=3000"`
	IsService   bool
}

func (e NewMessageEvent) Validate() error {
	return validator.Validator.Struct(e)
}

type MessageSentEvent struct {
	event `gonstructor:"-"`

	EventID   types.EventID   `validate:"required"`
	RequestID types.RequestID `validate:"required"`
	MessageID types.MessageID `validate:"required"`
}

func (e MessageSentEvent) Validate() error {
	return validator.Validator.Struct(e)
}

type MessageBlockedEvent struct {
	event `gonstructor:"-"`

	EventID   types.EventID   `validate:"required"`
	RequestID types.RequestID `validate:"required"`
	MessageID types.MessageID `validate:"required"`
}

func (e MessageBlockedEvent) Validate() error {
	return validator.Validator.Struct(e)
}

type NewChatEvent struct {
	event `gonstructor:"-"`

	EventID            types.EventID   `validate:"required"`
	RequestID          types.RequestID `validate:"required"`
	ChatID             types.ChatID    `validate:"required"`
	ClientID           types.UserID    `validate:"required"`
	CanTakeMoreProblem bool
}

func (e NewChatEvent) Validate() error {
	return validator.Validator.Struct(e)
}

type ChatClosedEvent struct {
	event `gonstructor:"-"`

	EventID            types.EventID   `validate:"required"`
	RequestID          types.RequestID `validate:"required"`
	ChatID             types.ChatID    `validate:"required"`
	CanTakeMoreProblem bool
}

func (e ChatClosedEvent) Validate() error {
	return validator.Validator.Struct(e)
}
