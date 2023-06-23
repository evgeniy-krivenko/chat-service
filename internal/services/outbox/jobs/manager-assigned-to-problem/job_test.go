package managerassignedtoproblemjob_test

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
	managerassignedtoproblemjob "github.com/evgeniy-krivenko/chat-service/internal/services/outbox/jobs/manager-assigned-to-problem"
	managerassignedtoproblemjobmocks "github.com/evgeniy-krivenko/chat-service/internal/services/outbox/jobs/manager-assigned-to-problem/mocks"
	"github.com/evgeniy-krivenko/chat-service/internal/testingh"
	"github.com/evgeniy-krivenko/chat-service/internal/types"
)

type JobSuite struct {
	testingh.ContextSuite

	ctrl        *gomock.Controller
	chatRepo    *managerassignedtoproblemjobmocks.MockchatRepo
	msgRepo     *managerassignedtoproblemjobmocks.MockmessageRepo
	mngrLoadSvc *managerassignedtoproblemjobmocks.MockmanagerLoadService
	eventStream *managerassignedtoproblemjobmocks.MockeventStream

	msgID     types.MessageID
	managerID types.UserID
	problemID types.ProblemID
	chatID    types.ChatID
	clientID  types.UserID
	reqID     types.RequestID

	message messagesrepo.Message
	chat    chatsrepo.Chat

	job *managerassignedtoproblemjob.Job
}

func TestJobSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(JobSuite))
}

func (j *JobSuite) SetupTest() {
	j.ContextSuite.SetupTest()

	j.ctrl = gomock.NewController(j.T())
	j.chatRepo = managerassignedtoproblemjobmocks.NewMockchatRepo(j.ctrl)
	j.msgRepo = managerassignedtoproblemjobmocks.NewMockmessageRepo(j.ctrl)
	j.mngrLoadSvc = managerassignedtoproblemjobmocks.NewMockmanagerLoadService(j.ctrl)
	j.eventStream = managerassignedtoproblemjobmocks.NewMockeventStream(j.ctrl)

	j.msgID = types.NewMessageID()
	j.managerID = types.NewUserID()
	j.problemID = types.NewProblemID()
	j.chatID = types.NewChatID()
	j.clientID = types.NewUserID()
	j.reqID = types.NewRequestID()
	createdAt := time.Now()

	j.message = messagesrepo.Message{
		ID:                  j.msgID,
		ChatID:              j.chatID,
		Body:                "go is awesome",
		CreatedAt:           createdAt,
		AuthorID:            types.UserIDNil,
		InitialRequestID:    j.reqID,
		IsVisibleForClient:  true,
		IsVisibleForManager: false,
		IsBlocked:           false,
		IsService:           true,
	}

	j.chat = chatsrepo.Chat{
		ID:        j.chatID,
		ClientID:  j.clientID,
		CreatedAt: createdAt,
	}

	var err error
	j.job, err = managerassignedtoproblemjob.New(managerassignedtoproblemjob.NewOptions(
		j.chatRepo,
		j.msgRepo,
		j.mngrLoadSvc,
		j.eventStream,
	))
	j.Require().NoError(err)
}

func (j *JobSuite) TearDown() {
	j.ContextSuite.TearDownTest()
	j.ctrl.Finish()
}

func (j *JobSuite) TestHandle_Success() {
	// Arrange.
	j.msgRepo.EXPECT().GetMessageByID(gomock.Any(), j.msgID).Return(&j.message, nil)
	j.chatRepo.EXPECT().GetChatByID(gomock.Any(), j.message.ChatID).Return(&j.chat, nil)
	j.mngrLoadSvc.EXPECT().CanManagerTakeProblem(gomock.Any(), j.managerID).Return(true, nil)
	j.eventStream.EXPECT().Publish(gomock.Any(), j.chat.ClientID, newMessageEventMatcher(eventstream.NewNewMessageEvent(
		types.NewEventID(),
		j.reqID,
		j.message.ChatID,
		j.message.ID,
		types.UserIDNil,
		j.message.CreatedAt,
		j.message.Body,
		"",
		j.message.IsService,
	))).Return(nil)
	j.eventStream.EXPECT().Publish(gomock.Any(), j.managerID, newChatEventMatcher(eventstream.NewNewChatEvent(
		types.NewEventID(),
		j.reqID,
		j.message.ChatID,
		j.clientID,
		true,
	))).Return(nil)

	// Action & assert
	payload, err := managerassignedtoproblemjob.MarshalPayload(j.msgID, j.managerID, j.problemID, j.reqID)
	j.Require().NoError(err)

	err = j.job.Handle(j.Ctx, payload)
	j.Require().NoError(err)
}

func (j *JobSuite) TestHandle_Error() {
	err := errors.New("unexpected")

	j.Run("get message error", func() {
		// Arrange.
		j.msgRepo.EXPECT().GetMessageByID(gomock.Any(), j.msgID).Return(nil, err)

		// Action & assert
		payload, err := managerassignedtoproblemjob.MarshalPayload(j.msgID, j.managerID, j.problemID, j.reqID)
		j.Require().NoError(err)

		err = j.job.Handle(j.Ctx, payload)
		j.Require().Error(err)
	})

	j.Run("get chat error", func() {
		// Arrange.
		j.msgRepo.EXPECT().GetMessageByID(gomock.Any(), j.msgID).Return(&j.message, nil)
		j.chatRepo.EXPECT().GetChatByID(gomock.Any(), j.message.ChatID).Return(nil, err)

		// Action & assert
		payload, err := managerassignedtoproblemjob.MarshalPayload(j.msgID, j.managerID, j.problemID, j.reqID)
		j.Require().NoError(err)

		err = j.job.Handle(j.Ctx, payload)
		j.Require().Error(err)
	})

	j.Run("get can manager take problem error", func() {
		// Arrange.
		j.msgRepo.EXPECT().GetMessageByID(gomock.Any(), j.msgID).Return(&j.message, nil)
		j.chatRepo.EXPECT().GetChatByID(gomock.Any(), j.message.ChatID).Return(&j.chat, nil)
		j.mngrLoadSvc.EXPECT().CanManagerTakeProblem(gomock.Any(), j.managerID).Return(false, err)

		// Action & assert
		payload, err := managerassignedtoproblemjob.MarshalPayload(j.msgID, j.managerID, j.problemID, j.reqID)
		j.Require().NoError(err)

		err = j.job.Handle(j.Ctx, payload)
		j.Require().Error(err)
	})

	j.Run("publish client event error", func() {
		// Arrange.
		j.msgRepo.EXPECT().GetMessageByID(gomock.Any(), j.msgID).Return(&j.message, nil)
		j.chatRepo.EXPECT().GetChatByID(gomock.Any(), j.message.ChatID).Return(&j.chat, nil)
		j.mngrLoadSvc.EXPECT().CanManagerTakeProblem(gomock.Any(), j.managerID).Return(false, nil)
		j.eventStream.EXPECT().Publish(gomock.Any(), j.chat.ClientID, newMessageEventMatcher(eventstream.NewNewMessageEvent(
			types.NewEventID(),
			j.reqID,
			j.message.ChatID,
			j.message.ID,
			types.UserIDNil,
			j.message.CreatedAt,
			j.message.Body,
			"",
			j.message.IsService,
		))).Return(err)

		// Action & assert
		payload, err := managerassignedtoproblemjob.MarshalPayload(j.msgID, j.managerID, j.problemID, j.reqID)
		j.Require().NoError(err)

		err = j.job.Handle(j.Ctx, payload)
		j.Require().Error(err)
	})

	j.Run("publish manager event error", func() {
		// Arrange.
		j.msgRepo.EXPECT().GetMessageByID(gomock.Any(), j.msgID).Return(&j.message, nil)
		j.chatRepo.EXPECT().GetChatByID(gomock.Any(), j.message.ChatID).Return(&j.chat, nil)
		j.mngrLoadSvc.EXPECT().CanManagerTakeProblem(gomock.Any(), j.managerID).Return(false, nil)
		j.eventStream.EXPECT().Publish(gomock.Any(), j.chat.ClientID, newMessageEventMatcher(eventstream.NewNewMessageEvent(
			types.NewEventID(),
			j.reqID,
			j.message.ChatID,
			j.message.ID,
			types.UserIDNil,
			j.message.CreatedAt,
			j.message.Body,
			"",
			j.message.IsService,
		))).Return(nil)
		j.eventStream.EXPECT().Publish(gomock.Any(), j.managerID, newChatEventMatcher(eventstream.NewNewChatEvent(
			types.NewEventID(),
			j.reqID,
			j.message.ChatID,
			j.clientID,
			false,
		))).Return(err)

		// Action & assert
		payload, err := managerassignedtoproblemjob.MarshalPayload(j.msgID, j.managerID, j.problemID, j.reqID)
		j.Require().NoError(err)

		err = j.job.Handle(j.Ctx, payload)
		j.Require().Error(err)
	})
}

type eqNewMessageEventParamsMatcher struct {
	arg *eventstream.NewMessageEvent
}

func newMessageEventMatcher(ev *eventstream.NewMessageEvent) gomock.Matcher {
	return &eqNewMessageEventParamsMatcher{arg: ev}
}

func (e *eqNewMessageEventParamsMatcher) Matches(x interface{}) bool {
	ev, ok := x.(*eventstream.NewMessageEvent)
	if !ok {
		return false
	}

	switch {
	case !e.arg.RequestID.Matches(ev.RequestID):
		return false
	case !e.arg.AuthorID.Matches(ev.AuthorID):
		return false
	case !e.arg.ChatID.Matches(ev.ChatID):
		return false
	case !e.arg.MessageID.Matches(ev.MessageID):
		return false
	case e.arg.MessageBody != ev.MessageBody:
		return false
	case e.arg.IsService != ev.IsService:
		return false
	case e.arg.CreatedAt.String() != ev.CreatedAt.String():
		return false
	}

	return true
}

func (e *eqNewMessageEventParamsMatcher) String() string {
	return fmt.Sprintf("%v", e.arg)
}

type eqNewChatEventParamsMatcher struct {
	arg *eventstream.NewChatEvent
}

func newChatEventMatcher(ev *eventstream.NewChatEvent) gomock.Matcher {
	return &eqNewChatEventParamsMatcher{arg: ev}
}

func (e *eqNewChatEventParamsMatcher) Matches(x interface{}) bool {
	ev, ok := x.(*eventstream.NewChatEvent)
	if !ok {
		return false
	}

	switch {
	case !e.arg.RequestID.Matches(ev.RequestID):
		return false
	case !e.arg.ChatID.Matches(ev.ChatID):
		return false
	case !e.arg.ClientID.Matches(ev.ClientID):
		return false
	case e.arg.CanTakeMoreProblem != ev.CanTakeMoreProblem:
		return false
	}

	return true
}

func (e *eqNewChatEventParamsMatcher) String() string {
	return fmt.Sprintf("%v", e.arg)
}
