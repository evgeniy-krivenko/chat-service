package eventstream_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	eventstream "github.com/evgeniy-krivenko/chat-service/internal/services/event-stream"
	"github.com/evgeniy-krivenko/chat-service/internal/types"
)

func TestNewMessageEvent_Validate(t *testing.T) {
	ev := eventstream.NewNewMessageEvent(
		types.MustParse[types.EventID]("d0ffbd36-bc30-11ed-8286-461e464ebed8"),
		types.MustParse[types.RequestID]("cee5f290-bc30-11ed-b7fe-461e464ebed8"),
		types.MustParse[types.ChatID]("31b4dc06-bc31-11ed-93cc-461e464ebed8"),
		types.MustParse[types.MessageID]("cb36a888-bc30-11ed-b843-461e464ebed8"),
		types.UserIDNil,
		time.Unix(1, 1).UTC(),
		"Manager will coming soon",
		"",
		true,
	)

	err := ev.Validate()
	assert.NoError(t, err)
}

func TestNewMessageSentEvent(t *testing.T) {
	ev := eventstream.NewMessageSentEvent(
		types.MustParse[types.EventID]("d0ffbd36-bc30-11ed-8286-461e464ebed8"),
		types.MustParse[types.RequestID]("cee5f290-bc30-11ed-b7fe-461e464ebed8"),
		types.MustParse[types.MessageID]("cb36a888-bc30-11ed-b843-461e464ebed8"),
	)

	err := ev.Validate()
	assert.NoError(t, err)
}

func TestNewMessageBlockedEvent(t *testing.T) {
	ev := eventstream.NewMessageBlockedEvent(
		types.MustParse[types.EventID]("d0ffbd36-bc30-11ed-8286-461e464ebed8"),
		types.MustParse[types.RequestID]("cee5f290-bc30-11ed-b7fe-461e464ebed8"),
		types.MustParse[types.MessageID]("cb36a888-bc30-11ed-b843-461e464ebed8"),
	)

	err := ev.Validate()
	assert.NoError(t, err)
}

func TestNewNewChatEvent(t *testing.T) {
	ev := eventstream.NewNewChatEvent(
		types.MustParse[types.EventID]("d0ffbd36-bc30-11ed-8286-461e464ebed8"),
		types.MustParse[types.RequestID]("cee5f290-bc30-11ed-b7fe-461e464ebed8"),
		types.MustParse[types.ChatID]("31b4dc06-bc31-11ed-93cc-461e464ebed8"),
		types.MustParse[types.UserID]("31b4dc06-bc31-11ed-93cc-461e464ebed8"),
		true,
	)

	err := ev.Validate()
	assert.NoError(t, err)
}
