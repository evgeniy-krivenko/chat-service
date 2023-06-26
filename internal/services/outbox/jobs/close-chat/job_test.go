package closechatjob_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	chatsrepo "github.com/evgeniy-krivenko/chat-service/internal/repositories/chats"
	messagesrepo "github.com/evgeniy-krivenko/chat-service/internal/repositories/messages"
	eventstream "github.com/evgeniy-krivenko/chat-service/internal/services/event-stream"
	closechatjob "github.com/evgeniy-krivenko/chat-service/internal/services/outbox/jobs/close-chat"
	closechatjobmocks "github.com/evgeniy-krivenko/chat-service/internal/services/outbox/jobs/close-chat/mocks"
	"github.com/evgeniy-krivenko/chat-service/internal/testingh"
	"github.com/evgeniy-krivenko/chat-service/internal/types"
)

type JobHandleSuite struct {
	testingh.ContextSuite

	ctrl        *gomock.Controller
	chatsRepo   *closechatjobmocks.MockchatsRepository
	msgRepo     *closechatjobmocks.MockmessageRepository
	eventStream *closechatjobmocks.MockeventStream
	mngrLoad    *closechatjobmocks.MockmanagerLoadService
	job         *closechatjob.Job

	reqID     types.RequestID
	chatID    types.ChatID
	managerID types.UserID
	clientID  types.UserID
	msgID     types.MessageID
}

func TestJobHandleSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(JobHandleSuite))
}

func (s *JobHandleSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.chatsRepo = closechatjobmocks.NewMockchatsRepository(s.ctrl)
	s.msgRepo = closechatjobmocks.NewMockmessageRepository(s.ctrl)
	s.eventStream = closechatjobmocks.NewMockeventStream(s.ctrl)
	s.mngrLoad = closechatjobmocks.NewMockmanagerLoadService(s.ctrl)

	s.reqID = types.NewRequestID()
	s.chatID = types.NewChatID()
	s.managerID = types.NewUserID()
	s.clientID = types.NewUserID()
	s.msgID = types.NewMessageID()

	var err error
	s.job, err = closechatjob.New(closechatjob.NewOptions(s.mngrLoad, s.eventStream, s.chatsRepo, s.msgRepo))
	s.Require().NoError(err)

	s.ContextSuite.SetupTest()
}

func (s *JobHandleSuite) TearDownTest() {
	s.ctrl.Finish()

	s.ContextSuite.TearDownTest()
}

func (s *JobHandleSuite) TestHandle_Success() {
	// Arrange.
	msg := messagesrepo.Message{
		ID:        s.msgID,
		ChatID:    s.chatID,
		Body:      "some service msg",
		IsService: true,
		CreatedAt: time.Now(),
	}

	payload, err := closechatjob.MarshalPayload(s.reqID, s.managerID, s.chatID, s.msgID)
	s.Require().NoError(err)
	s.mngrLoad.EXPECT().CanManagerTakeProblem(gomock.Any(), s.managerID).Return(true, nil)
	s.msgRepo.EXPECT().GetMessageByID(gomock.Any(), s.msgID).Return(&msg, nil)
	s.chatsRepo.EXPECT().GetChatByID(gomock.Any(), s.chatID).Return(&chatsrepo.Chat{
		ClientID: s.clientID,
	}, nil)

	s.eventStream.EXPECT().Publish(gomock.Any(), s.clientID, newMessageEventMatcher{
		NewMessageEvent: &eventstream.NewMessageEvent{
			EventID:     types.EventIDNil,
			RequestID:   s.reqID,
			ChatID:      msg.ChatID,
			MessageID:   msg.ID,
			AuthorID:    types.UserIDNil, // no possible to check
			CreatedAt:   msg.CreatedAt,
			MessageBody: "some service msg",
			IsService:   msg.IsService,
		},
	})

	s.eventStream.EXPECT().Publish(gomock.Any(), s.managerID, chatClosedEventMatcher{
		ChatClosedEvent: &eventstream.ChatClosedEvent{
			EventID:            types.EventIDNil,
			RequestID:          s.reqID,
			ChatID:             s.chatID,
			CanTakeMoreProblem: true,
		},
	})

	// Action.
	err = s.job.Handle(s.Ctx, payload)

	// Assert.
	s.Require().NoError(err)
}

func (s *JobHandleSuite) TestHandle_UnmarshalError() {
	// Arrange.
	wrongPayload := `{"RequestID":`

	// Action.
	err := s.job.Handle(s.Ctx, wrongPayload)

	// Assert.
	s.Require().Error(err)
}

func (s *JobHandleSuite) TestHandle_CanManagerTakeProblemError() {
	// Arrange.
	payload, err := closechatjob.MarshalPayload(s.reqID, s.managerID, s.chatID, s.msgID)
	s.Require().NoError(err)
	s.mngrLoad.EXPECT().CanManagerTakeProblem(gomock.Any(), s.managerID).
		Return(false, errors.New("unexpected"))

	// Action.
	err = s.job.Handle(s.Ctx, payload)

	// Assert.
	s.Require().Error(err)
}

func (s *JobHandleSuite) TestHandle_GetMessageByIDError() {
	// Arrange.
	payload, err := closechatjob.MarshalPayload(s.reqID, s.managerID, s.chatID, s.msgID)
	s.Require().NoError(err)
	s.mngrLoad.EXPECT().CanManagerTakeProblem(gomock.Any(), s.managerID).Return(true, nil)
	s.msgRepo.EXPECT().GetMessageByID(gomock.Any(), s.msgID).Return(nil, errors.New("unexpected"))

	// Action.
	err = s.job.Handle(s.Ctx, payload)

	// Assert.
	s.Require().Error(err)
}

func (s *JobHandleSuite) TestHandle_GetChatByIDError() {
	// Arrange.
	msg := messagesrepo.Message{
		ID:        s.msgID,
		ChatID:    s.chatID,
		Body:      "some service msg",
		IsService: true,
		CreatedAt: time.Now(),
	}

	payload, err := closechatjob.MarshalPayload(s.reqID, s.managerID, s.chatID, s.msgID)
	s.Require().NoError(err)
	s.mngrLoad.EXPECT().CanManagerTakeProblem(gomock.Any(), s.managerID).Return(true, nil)
	s.msgRepo.EXPECT().GetMessageByID(gomock.Any(), s.msgID).Return(&msg, nil)
	s.chatsRepo.EXPECT().GetChatByID(gomock.Any(), s.chatID).Return(nil, errors.New("unexpected"))

	// Action.
	err = s.job.Handle(s.Ctx, payload)

	// Assert.
	s.Require().Error(err)
}

func (s *JobHandleSuite) TestHandle_ClientPublishError() {
	// Arrange.
	msg := messagesrepo.Message{
		ID:        s.msgID,
		ChatID:    s.chatID,
		Body:      "some service msg",
		IsService: true,
		CreatedAt: time.Now(),
	}
	payload, err := closechatjob.MarshalPayload(s.reqID, s.managerID, s.chatID, s.msgID)
	s.Require().NoError(err)
	s.mngrLoad.EXPECT().CanManagerTakeProblem(gomock.Any(), s.managerID).Return(true, nil)
	s.msgRepo.EXPECT().GetMessageByID(gomock.Any(), s.msgID).Return(&msg, nil)
	s.chatsRepo.EXPECT().GetChatByID(gomock.Any(), s.chatID).Return(&chatsrepo.Chat{
		ClientID: s.clientID,
	}, nil)

	s.eventStream.EXPECT().Publish(gomock.Any(), s.clientID, newMessageEventMatcher{
		NewMessageEvent: &eventstream.NewMessageEvent{
			EventID:     types.EventIDNil,
			RequestID:   s.reqID,
			ChatID:      msg.ChatID,
			MessageID:   msg.ID,
			AuthorID:    types.UserIDNil, // no possible to check
			CreatedAt:   msg.CreatedAt,
			MessageBody: "some service msg",
			IsService:   msg.IsService,
		},
	}).Return(errors.New("unexpected"))

	s.eventStream.EXPECT().Publish(gomock.Any(), s.managerID, chatClosedEventMatcher{
		ChatClosedEvent: &eventstream.ChatClosedEvent{
			EventID:            types.EventIDNil,
			RequestID:          s.reqID,
			ChatID:             s.chatID,
			CanTakeMoreProblem: true,
		},
	}).AnyTimes()

	// Action.
	err = s.job.Handle(s.Ctx, payload)

	// Assert.
	s.Require().Error(err)
}

func (s *JobHandleSuite) TestHandle_ManagerPublishError() {
	// Arrange.
	msg := messagesrepo.Message{
		ID:        s.msgID,
		ChatID:    s.chatID,
		Body:      "some service msg",
		IsService: true,
		CreatedAt: time.Now(),
	}

	payload, err := closechatjob.MarshalPayload(s.reqID, s.managerID, s.chatID, s.msgID)
	s.Require().NoError(err)
	s.mngrLoad.EXPECT().CanManagerTakeProblem(gomock.Any(), s.managerID).Return(true, nil)
	s.msgRepo.EXPECT().GetMessageByID(gomock.Any(), s.msgID).Return(&msg, nil)
	s.chatsRepo.EXPECT().GetChatByID(gomock.Any(), s.chatID).Return(&chatsrepo.Chat{
		ClientID: s.clientID,
	}, nil)

	s.eventStream.EXPECT().Publish(gomock.Any(), s.clientID, newMessageEventMatcher{
		NewMessageEvent: &eventstream.NewMessageEvent{
			EventID:     types.EventIDNil,
			RequestID:   s.reqID,
			ChatID:      msg.ChatID,
			MessageID:   msg.ID,
			AuthorID:    types.UserIDNil, // no possible to check
			CreatedAt:   msg.CreatedAt,
			MessageBody: "some service msg",
			IsService:   msg.IsService,
		},
	}).AnyTimes()

	s.eventStream.EXPECT().Publish(gomock.Any(), s.managerID, chatClosedEventMatcher{
		ChatClosedEvent: &eventstream.ChatClosedEvent{
			EventID:            types.EventIDNil,
			RequestID:          s.reqID,
			ChatID:             s.chatID,
			CanTakeMoreProblem: true,
		},
	}).Return(errors.New("unexpected"))

	// Action.
	err = s.job.Handle(s.Ctx, payload)

	// Assert.
	s.Require().Error(err)
}

var _ gomock.Matcher = newMessageEventMatcher{}

type newMessageEventMatcher struct {
	*eventstream.NewMessageEvent
}

func (m newMessageEventMatcher) Matches(x interface{}) bool {
	envelope, ok := x.(eventstream.Event)
	if !ok {
		return false
	}

	ev, ok := envelope.(*eventstream.NewMessageEvent)
	if !ok {
		return false
	}

	return !ev.EventID.IsZero() &&
		ev.RequestID == m.RequestID &&
		ev.ChatID == m.ChatID &&
		ev.MessageID == m.MessageID &&
		ev.CreatedAt.Equal(m.CreatedAt) &&
		ev.MessageBody == m.MessageBody &&
		ev.IsService == m.IsService
}

func (m newMessageEventMatcher) String() string {
	return fmt.Sprintf("%v", m.NewMessageEvent)
}

var _ gomock.Matcher = chatClosedEventMatcher{}

type chatClosedEventMatcher struct {
	*eventstream.ChatClosedEvent
}

func (m chatClosedEventMatcher) Matches(x interface{}) bool {
	envelope, ok := x.(eventstream.Event)
	if !ok {
		return false
	}

	ev, ok := envelope.(*eventstream.ChatClosedEvent)
	if !ok {
		return false
	}

	return !ev.EventID.IsZero() &&
		ev.RequestID == m.RequestID &&
		ev.ChatID == m.ChatID &&
		ev.CanTakeMoreProblem == m.CanTakeMoreProblem
}

func (m chatClosedEventMatcher) String() string {
	return fmt.Sprintf("%v", m.ChatClosedEvent)
}
