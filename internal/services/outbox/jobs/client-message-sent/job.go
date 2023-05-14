package clientmessagesentjob

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	messagesrepo "github.com/evgeniy-krivenko/chat-service/internal/repositories/messages"
	eventstream "github.com/evgeniy-krivenko/chat-service/internal/services/event-stream"
	"github.com/evgeniy-krivenko/chat-service/internal/services/outbox"
	msgjobpayload "github.com/evgeniy-krivenko/chat-service/internal/services/outbox/msg-job-payload"
	"github.com/evgeniy-krivenko/chat-service/internal/types"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/job_mock.gen.go -package=clientmessagesentjobmocks

const Name = "client-message-sent"

type messageRepository interface {
	GetMessageByID(ctx context.Context, msgID types.MessageID) (*messagesrepo.Message, error)
}

type eventStream interface {
	Publish(ctx context.Context, userID types.UserID, event eventstream.Event) error
}

//go:generate options-gen -out-filename=job_options.gen.go -from-struct=Options
type Options struct {
	msgRepo     messageRepository `option:"mandatory" validate:"required"`
	eventStream eventStream       `option:"mandatory" validate:"required"`
}

type Job struct {
	outbox.DefaultJob
	Options
	lg *zap.Logger
}

func New(opts Options) (*Job, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("validate job %v options: %v", Name, err)
	}

	return &Job{
		Options: opts,
		lg:      zap.L().Named(Name),
	}, nil
}

func (j *Job) Name() string {
	return Name
}

func (j *Job) Handle(ctx context.Context, payload string) error {
	msgID, err := msgjobpayload.UnmarshalPayload(payload)
	if err != nil {
		j.lg.Error("unmarshal payload", zap.Error(err))
		return fmt.Errorf("unmarshal payload: %v", err)
	}

	msg, err := j.msgRepo.GetMessageByID(ctx, msgID)
	if err != nil {
		j.lg.Error("get msg from repo", zap.Error(err))
		return fmt.Errorf("get msg while hadle job %v: %v", Name, err)
	}

	if err := j.eventStream.Publish(ctx, msg.AuthorID, eventstream.NewMessageSentEvent(
		types.NewEventID(),
		msg.InitialRequestID,
		msg.ID,
	)); err != nil {
		return fmt.Errorf("publish msg sent to event stream: %v", err)
	}
	return nil
}
