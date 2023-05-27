package clientmessagesentjob_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	messagesrepo "github.com/evgeniy-krivenko/chat-service/internal/repositories/messages"
	eventstream "github.com/evgeniy-krivenko/chat-service/internal/services/event-stream"
	clientmessagesentjob "github.com/evgeniy-krivenko/chat-service/internal/services/outbox/jobs/client-message-sent"
	clientmessagesentjobmocks "github.com/evgeniy-krivenko/chat-service/internal/services/outbox/jobs/client-message-sent/mocks"
	"github.com/evgeniy-krivenko/chat-service/internal/services/outbox/jobs/payload/simpleid"
	"github.com/evgeniy-krivenko/chat-service/internal/testingh"
	"github.com/evgeniy-krivenko/chat-service/internal/types"
)

type JobSuite struct {
	testingh.ContextSuite

	ctrl        *gomock.Controller
	msgRepo     *clientmessagesentjobmocks.MockmessageRepository
	eventStream *clientmessagesentjobmocks.MockeventStream
	job         *clientmessagesentjob.Job
	msg         *messagesrepo.Message
}

func TestJobSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(JobSuite))
}

func (j *JobSuite) SetupTest() {
	j.ContextSuite.SetupTest()

	j.ctrl = gomock.NewController(j.T())
	j.msgRepo = clientmessagesentjobmocks.NewMockmessageRepository(j.ctrl)
	j.eventStream = clientmessagesentjobmocks.NewMockeventStream(j.ctrl)
	j.msg = &messagesrepo.Message{
		ID:               types.NewMessageID(),
		InitialRequestID: types.NewRequestID(),
		AuthorID:         types.NewUserID(),
	}

	var err error
	j.job, err = clientmessagesentjob.New(clientmessagesentjob.NewOptions(j.msgRepo, j.eventStream))
	j.Require().NoError(err)
}

func (j *JobSuite) TearDown() {
	j.ContextSuite.TearDownTest()
	j.ctrl.Finish()
}

func (j *JobSuite) TestHandle_Success() {
	// Arrange.
	j.msgRepo.EXPECT().GetMessageByID(gomock.Any(), j.msg.ID).Return(j.msg, nil)
	j.eventStream.EXPECT().Publish(gomock.Any(), j.msg.AuthorID, newMessageEventMatcher(eventstream.NewMessageSentEvent(
		types.NewEventID(),
		j.msg.InitialRequestID,
		j.msg.ID,
	))).Return(nil)

	// Action & assert.
	payload, err := simpleid.Marshal(j.msg.ID)
	j.Require().NoError(err)

	err = j.job.Handle(j.Ctx, payload)
	j.Require().NoError(err)
}

func (j *JobSuite) TestHandle_ErrorMsgRepo() {
	// Arrange.
	j.msgRepo.EXPECT().GetMessageByID(gomock.Any(), j.msg.ID).
		Return(nil, errors.New("unexpected"))

	// Action & assert.
	payload, err := simpleid.Marshal(j.msg.ID)
	j.Require().NoError(err)

	err = j.job.Handle(j.Ctx, payload)
	j.Require().Error(err)
}

func (j *JobSuite) TestHandle_ErrorPublish() {
	// Arrange.
	j.msgRepo.EXPECT().GetMessageByID(gomock.Any(), j.msg.ID).Return(j.msg, nil)
	j.eventStream.EXPECT().Publish(gomock.Any(), j.msg.AuthorID, newMessageEventMatcher(eventstream.NewMessageSentEvent(
		types.NewEventID(),
		j.msg.InitialRequestID,
		j.msg.ID,
	))).Return(errors.New("unexpected"))

	// Action & assert.
	payload, err := simpleid.Marshal(j.msg.ID)
	j.Require().NoError(err)

	err = j.job.Handle(j.Ctx, payload)
	j.Require().Error(err)
}

type eqNewMessageEventParamsMatcher struct {
	arg *eventstream.MessageSentEvent
}

func newMessageEventMatcher(ev *eventstream.MessageSentEvent) gomock.Matcher {
	return &eqNewMessageEventParamsMatcher{arg: ev}
}

func (e *eqNewMessageEventParamsMatcher) Matches(x interface{}) bool {
	ev, ok := x.(*eventstream.MessageSentEvent)
	if !ok {
		return false
	}

	switch {
	case e.arg.RequestID.String() != ev.RequestID.String():
		return false
	case !e.arg.MessageID.Matches(ev.MessageID):
		return false
	}

	return true
}

func (e *eqNewMessageEventParamsMatcher) String() string {
	return fmt.Sprintf("%v", e.arg)
}
