package managerevents

import (
	"fmt"

	eventstream "github.com/evgeniy-krivenko/chat-service/internal/services/event-stream"
	websocketstream "github.com/evgeniy-krivenko/chat-service/internal/websocket-stream"
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
	default:
		return nil, fmt.Errorf("unknown manager event: %v (%T)", e, e)
	}
	if err != nil {
		return nil, err
	}

	return event, nil
}
