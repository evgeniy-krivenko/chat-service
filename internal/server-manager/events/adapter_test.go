package managerevents_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	managerevents "github.com/evgeniy-krivenko/chat-service/internal/server-manager/events"
	eventstream "github.com/evgeniy-krivenko/chat-service/internal/services/event-stream"
	"github.com/evgeniy-krivenko/chat-service/internal/types"
)

func TestAdapter_Adapt(t *testing.T) {
	cases := []struct {
		name    string
		ev      eventstream.Event
		expJSON string
	}{
		{
			name: "chat event",
			ev: eventstream.NewNewChatEvent(
				types.MustParse[types.EventID]("d0ffbd36-bc30-11ed-8286-461e464ebed8"),
				types.MustParse[types.RequestID]("cee5f290-bc30-11ed-b7fe-461e464ebed8"),
				types.MustParse[types.ChatID]("cb36a888-bc30-11ed-b843-461e464ebed8"),
				types.MustParse[types.UserID]("cee5f290-bc30-11ed-b7fe-461e464ebed8"),
				false,
			),
			expJSON: `{
				"eventId": "d0ffbd36-bc30-11ed-8286-461e464ebed8",
				"requestId": "cee5f290-bc30-11ed-b7fe-461e464ebed8",
				"chatId": "cb36a888-bc30-11ed-b843-461e464ebed8",
				"clientId": "cee5f290-bc30-11ed-b7fe-461e464ebed8",
				"eventType":"NewChatEvent",
				"canTakeMoreProblems": false
			}`,
		},
	}

	for _, tt := range cases {
		t.Run(
			tt.name,
			func(t *testing.T) {
				adapted, err := managerevents.Adapter{}.Adapt(tt.ev)
				require.NoError(t, err)

				raw, err := json.Marshal(adapted)
				require.NoError(t, err)
				assert.JSONEq(t, tt.expJSON, string(raw))
			},
		)
	}
}
