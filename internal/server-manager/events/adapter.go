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

	switch e := ev.(type) {
	case *eventstream.NewChatEvent:
		return &NewChatEvent{
			EventID:             e.EventID,
			ClientID:            e.ClientID,
			CanTakeMoreProblems: e.CanTakeMoreProblem,
			RequestID:           e.RequestID,
			EventType:           EventTypeNewChatEvent,
		}, nil
	}

	return nil, nil
}
