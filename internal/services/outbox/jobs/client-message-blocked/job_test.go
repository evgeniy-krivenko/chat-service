package clientmessageblockedjob_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	messagesrepo "github.com/evgeniy-krivenko/chat-service/internal/repositories/messages"
	eventstream "github.com/evgeniy-krivenko/chat-service/internal/services/event-stream"
	clientmessageblockedjob "github.com/evgeniy-krivenko/chat-service/internal/services/outbox/jobs/client-message-blocked"
	clientmessageblockedeventmocks "github.com/evgeniy-krivenko/chat-service/internal/services/outbox/jobs/client-message-blocked/mocks"
	msgjobpayload "github.com/evgeniy-krivenko/chat-service/internal/services/outbox/msg-job-payload"
	"github.com/evgeniy-krivenko/chat-service/internal/testingh"
	"github.com/evgeniy-krivenko/chat-service/internal/types"
)

type JobSuite struct {
	testingh.ContextSuite

	ctrl        *gomock.Controller
	msgRepo     *clientmessageblockedeventmocks.MockmessageRepository
	eventStream *clientmessageblockedeventmocks.MockeventStream
	job         *clientmessageblockedjob.Job
	msg         *messagesrepo.Message
}

func TestJobSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(JobSuite))
}

func (j *JobSuite) SetupTest() {
	j.ContextSuite.SetupTest()

	j.ctrl = gomock.NewController(j.T())
	j.msgRepo = clientmessageblockedeventmocks.NewMockmessageRepository(j.ctrl)
	j.eventStream = clientmessageblockedeventmocks.NewMockeventStream(j.ctrl)
	j.msg = &messagesrepo.Message{
		ID:               types.NewMessageID(),
		InitialRequestID: types.NewRequestID(),
		AuthorID:         types.NewUserID(),
	}

	var err error
	j.job, err = clientmessageblockedjob.New(clientmessageblockedjob.NewOptions(j.msgRepo, j.eventStream))
	j.Require().NoError(err)
}

func (j *JobSuite) TearDown() {
	j.ContextSuite.TearDownTest()
	j.ctrl.Finish()
}

func (j *JobSuite) TestHandle_Success() {
	// Arrange.
	j.msgRepo.EXPECT().GetMessageByID(gomock.Any(), j.msg.ID).Return(j.msg, nil)
	j.eventStream.EXPECT().Publish(gomock.Any(), j.msg.AuthorID, newMessageBlockedEventMatcher(eventstream.NewMessageBlockedEvent(
		types.NewEventID(),
		j.msg.InitialRequestID,
		j.msg.ID,
	))).Return(nil)

	// Action & assert.
	payload, err := msgjobpayload.MarshalPayload(j.msg.ID)
	j.Require().NoError(err)

	err = j.job.Handle(j.Ctx, payload)
	j.Require().NoError(err)
}

func (j *JobSuite) TestHandle_ErrorMsgRepo() {
	// Arrange.
	j.msgRepo.EXPECT().GetMessageByID(gomock.Any(), j.msg.ID).
		Return(nil, errors.New("unexpected"))

	// Action & assert.
	payload, err := msgjobpayload.MarshalPayload(j.msg.ID)
	j.Require().NoError(err)

	err = j.job.Handle(j.Ctx, payload)
	j.Require().Error(err)
}

func (j *JobSuite) TestHandle_ErrorPublish() {
	// Arrange.
	j.msgRepo.EXPECT().GetMessageByID(gomock.Any(), j.msg.ID).Return(j.msg, nil)
	j.eventStream.EXPECT().Publish(gomock.Any(), j.msg.AuthorID, newMessageBlockedEventMatcher(eventstream.NewMessageBlockedEvent(
		types.NewEventID(),
		j.msg.InitialRequestID,
		j.msg.ID,
	))).Return(errors.New("unexpected"))

	// Action & assert.
	payload, err := msgjobpayload.MarshalPayload(j.msg.ID)
	j.Require().NoError(err)

	err = j.job.Handle(j.Ctx, payload)
	j.Require().Error(err)
}

type eqMessageBlockedEventParamsMatcher struct {
	arg *eventstream.MessageBlockedEvent
}

func newMessageBlockedEventMatcher(ev *eventstream.MessageBlockedEvent) gomock.Matcher {
	return &eqMessageBlockedEventParamsMatcher{arg: ev}
}

func (e *eqMessageBlockedEventParamsMatcher) Matches(x interface{}) bool {
	ev, ok := x.(*eventstream.MessageBlockedEvent)
	if !ok {
		return false
	}

	switch {
	case e.arg.RequestID.String() != ev.RequestID.String():
		return false
	case e.arg.MessageID.String() != ev.MessageID.String():
		return false
	}

	return true
}

func (e *eqMessageBlockedEventParamsMatcher) String() string {
	return fmt.Sprintf("%v", e.arg)
}
