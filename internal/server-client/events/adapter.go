package clientevents

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
		return nil, fmt.Errorf("validate while adapt: %v", err)
	}

	var event Event
	var err error

	switch e := ev.(type) {
	case *eventstream.NewMessageEvent:
		event.EventId = e.EventID
		event.RequestId = e.RequestID

		err = event.FromNewMessageEvent(NewMessageEvent{
			AuthorId:  pointer.PtrWithZeroAsNil(e.AuthorID),
			CreatedAt: e.CreatedAt,
			IsService: e.IsService,
			Body:      e.MessageBody,
			MessageId: e.MessageID,
		})
	case *eventstream.MessageSentEvent:
		event.EventId = e.EventID
		event.RequestId = e.RequestID

		err = event.FromMessageSentEvent(MessageSentEvent{
			MessageId: e.MessageID,
		})
	case *eventstream.MessageBlockedEvent:
		event.EventId = e.EventID
		event.RequestId = e.RequestID

		err = event.FromMessageBlockedEvent(MessageBlockedEvent{
			MessageId: e.MessageID,
		})
	default:
		return nil, fmt.Errorf("unknown client event: %v (%T)", e, e)
	}
	if err != nil {
		return nil, err
	}

	return event, nil
}
