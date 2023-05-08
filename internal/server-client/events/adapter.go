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

	switch e := ev.(type) {
	case *eventstream.NewMessageEvent:
		return &NewMessageEvent{
			EventID:   e.EventID,
			MessageID: e.MessageID,
			RequestID: e.RequestID,
			CreatedAt: e.CreatedAt,
			IsService: e.IsService,
			Body:      e.MessageBody,
			EventType: EventTypeNewMessageEvent,
			AuthorID:  pointer.PtrWithZeroAsNil(e.AuthorID),
		}, nil
	case *eventstream.MessageSentEvent:
		return &MessageSentEvent{
			EventID:   e.EventID,
			MessageID: e.MessageID,
			RequestID: e.RequestID,
			EventType: EventTypeMessageSentEvent,
		}, nil
	}
	return nil, nil
}
