package managerevents

import (
	"fmt"

	eventstream "github.com/evgeniy-krivenko/chat-service/internal/services/event-stream"
	websocketstream "github.com/evgeniy-krivenko/chat-service/internal/websocket-stream"
	"github.com/evgeniy-krivenko/chat-service/pkg/pointer"
)

var _ websocketstream.EventAdapter = Adapter{}

type Adapter struct{}

func (Adapter) Adapt(ev eventstream.Event) (any, error) {
	if err := ev.Validate(); err != nil {
		return nil, fmt.Errorf("validate while apapt event: %v", err)
	}

	var event Event
	var err error

	switch e := ev.(type) {
	case *eventstream.NewChatEvent:
		event.EventId = e.EventID
		event.RequestId = e.RequestID

		err = event.FromNewChatEvent(NewChatEvent{
			CanTakeMoreProblems: e.CanTakeMoreProblem,
			ClientId:            e.ClientID,
			ChatId:              e.ChatID,
			FirstName:           pointer.PtrWithZeroAsNil(e.FirstName),
			LastName:            pointer.PtrWithZeroAsNil(e.LastName),
		})
	case *eventstream.NewMessageEvent:
		event.EventId = e.EventID
		event.RequestId = e.RequestID

		err = event.FromNewMessageEvent(NewMessageEvent{
			AuthorId:  e.AuthorID,
			CreatedAt: e.CreatedAt,
			ChatId:    e.ChatID,
			Body:      e.MessageBody,
			MessageId: e.MessageID,
		})
	case *eventstream.ChatClosedEvent:
		event.EventId = e.EventID
		event.RequestId = e.RequestID

		err = event.FromChatClosedEvent(ChatClosedEvent{
			ChatId:              e.ChatID,
			CanTakeMoreProblems: e.CanTakeMoreProblem,
		})
	default:
		return nil, fmt.Errorf("unknown manager event: %v (%T)", e, e)
	}
	if err != nil {
		return nil, err
	}

	return event, nil
}
