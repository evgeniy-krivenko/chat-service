package sendmanagermessagejob_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	chatsrepo "github.com/evgeniy-krivenko/chat-service/internal/repositories/chats"
	messagesrepo "github.com/evgeniy-krivenko/chat-service/internal/repositories/messages"
	eventstream "github.com/evgeniy-krivenko/chat-service/internal/services/event-stream"
	msgproducer "github.com/evgeniy-krivenko/chat-service/internal/services/msg-producer"
	"github.com/evgeniy-krivenko/chat-service/internal/services/outbox/jobs/payload/simpleid"
	sendmanagermessagejob "github.com/evgeniy-krivenko/chat-service/internal/services/outbox/jobs/send-manager-message"
	sendmanagermessagejobmocks "github.com/evgeniy-krivenko/chat-service/internal/services/outbox/jobs/send-manager-message/mocks"
	"github.com/evgeniy-krivenko/chat-service/internal/types"
)

func TestJob_Handle(t *testing.T) {
	// Arrange.
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	msgProducer := sendmanagermessagejobmocks.NewMockmessageProducer(ctrl)
	msgRepo := sendmanagermessagejobmocks.NewMockmessageRepository(ctrl)
	eventStream := sendmanagermessagejobmocks.NewMockeventStream(ctrl)
	chatsRepo := sendmanagermessagejobmocks.NewMockchatRepository(ctrl)
	job, err := sendmanagermessagejob.New(sendmanagermessagejob.NewOptions(msgProducer, msgRepo, eventStream, chatsRepo))
	require.NoError(t, err)

	managerID := types.NewUserID()
	clientID := types.NewUserID()
	msgID := types.NewMessageID()
	chatID := types.NewChatID()
	reqID := types.NewRequestID()
	createdAt := time.Now()
	const body = "Hello!"

	msg := messagesrepo.Message{
		ID:                  msgID,
		ChatID:              chatID,
		AuthorID:            managerID,
		InitialRequestID:    reqID,
		Body:                body,
		CreatedAt:           createdAt,
		ManagerID:           types.UserIDNil,
		IsVisibleForClient:  true,
		IsVisibleForManager: true,
		IsBlocked:           false,
		IsService:           false,
	}
	chat := chatsrepo.Chat{
		ID:       chatID,
		ClientID: clientID,
	}

	msgRepo.EXPECT().GetMessageByID(gomock.Any(), msgID).Return(&msg, nil)

	chatsRepo.EXPECT().GetChatByID(gomock.Any(), msg.ChatID).Return(&chat, nil)

	msgProducer.EXPECT().ProduceMessage(gomock.Any(), msgproducer.Message{
		ID:         msgID,
		ChatID:     chatID,
		Body:       body,
		FromClient: false,
	}).Return(nil)

	eventStream.EXPECT().Publish(gomock.Any(), msg.AuthorID, newMessageEventMatcher(eventstream.NewNewMessageEvent(
		types.NewEventID(),
		reqID,
		chatID,
		msgID,
		managerID,
		createdAt,
		body,
		false,
	),
	)).Return(nil)

	eventStream.EXPECT().Publish(gomock.Any(), chat.ClientID, newMessageEventMatcher(eventstream.NewNewMessageEvent(
		types.NewEventID(),
		reqID,
		chatID,
		msgID,
		managerID,
		createdAt,
		body,
		false,
	),
	)).Return(nil)

	// Action & assert.
	payload, err := simpleid.Marshal(msgID)
	require.NoError(t, err)

	err = job.Handle(ctx, payload)
	require.NoError(t, err)
}

type eqNewMessageEventParamsMatcher struct {
	arg *eventstream.NewMessageEvent
}

func newMessageEventMatcher(ev *eventstream.NewMessageEvent) gomock.Matcher {
	return &eqNewMessageEventParamsMatcher{arg: ev}
}

func (e *eqNewMessageEventParamsMatcher) Matches(x interface{}) bool {
	envelope, ok := x.(eventstream.Event)
	if !ok {
		return false
	}

	ev, ok := envelope.(*eventstream.NewMessageEvent)
	if !ok {
		return false
	}

	return !ev.EventID.IsZero() &&
		ev.RequestID == e.arg.RequestID &&
		ev.ChatID == e.arg.ChatID &&
		ev.MessageID == e.arg.MessageID &&
		ev.AuthorID == e.arg.AuthorID &&
		ev.CreatedAt.Equal(e.arg.CreatedAt) &&
		ev.MessageBody == e.arg.MessageBody &&
		ev.IsService == e.arg.IsService
}

func (e *eqNewMessageEventParamsMatcher) String() string {
	return fmt.Sprintf("%v", e.arg)
}
