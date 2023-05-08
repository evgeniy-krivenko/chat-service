package eventstream_test

import (
	eventstream "github.com/evgeniy-krivenko/chat-service/internal/services/event-stream"
	"github.com/evgeniy-krivenko/chat-service/internal/types"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
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
		true,
	)

	err := ev.Validate()
	assert.NoError(t, err)
}
