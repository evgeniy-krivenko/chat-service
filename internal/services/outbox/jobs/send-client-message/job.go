package sendclientmessagejob

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	messagesrepo "github.com/evgeniy-krivenko/chat-service/internal/repositories/messages"
	eventstream "github.com/evgeniy-krivenko/chat-service/internal/services/event-stream"
	msgproducer "github.com/evgeniy-krivenko/chat-service/internal/services/msg-producer"
	"github.com/evgeniy-krivenko/chat-service/internal/services/outbox"
	"github.com/evgeniy-krivenko/chat-service/internal/services/outbox/jobs/payload/simpleid"
	"github.com/evgeniy-krivenko/chat-service/internal/types"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/job_mock.gen.go -package=sendclientmessagejobmocks

const Name = "send-client-message"

type messageProducer interface {
	ProduceMessage(ctx context.Context, message msgproducer.Message) error
}

type messageRepository interface {
	GetMessageByID(ctx context.Context, msgID types.MessageID) (*messagesrepo.Message, error)
}

type eventStream interface {
	Publish(ctx context.Context, userID types.UserID, event eventstream.Event) error
}

//go:generate options-gen -out-filename=job_options.gen.go -from-struct=Options
type Options struct {
	msgProducer messageProducer   `option:"mandatory" validate:"required"`
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
		return nil, fmt.Errorf("validate send client message job: %v", err)
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
	j.lg.Info("start processing", zap.String("payload", payload))

	msgID, err := simpleid.Unmarshal[types.MessageID](payload)
	if err != nil {
		j.lg.Warn("unmarshal payload", zap.Error(err))
		return fmt.Errorf("unmarshal payload: %v", err)
	}

	msg, err := j.msgRepo.GetMessageByID(ctx, msgID)
	if err != nil {
		j.lg.Warn("get msg from repo", zap.Error(err))
		return fmt.Errorf("get msg from repo: %v", err)
	}

	if err := j.msgProducer.ProduceMessage(ctx, msgproducer.Message{
		ID:         msg.ID,
		ChatID:     msg.ChatID,
		Body:       msg.Body,
		FromClient: true,
	}); err != nil {
		j.lg.Warn("produce message", zap.Error(err))
		return fmt.Errorf("produce msg: %v", err)
	}

	if err := j.eventStream.Publish(ctx, msg.AuthorID, eventstream.NewNewMessageEvent(
		types.NewEventID(),
		msg.InitialRequestID,
		msg.ChatID,
		msg.ID,
		msg.AuthorID,
		msg.CreatedAt,
		msg.Body,
		"",
		msg.IsService,
	)); err != nil {
		j.lg.Warn("publish message", zap.Stringer("message_id", msgID))
		return fmt.Errorf("publish NewMesaggeEvent to client stream: %v", err)
	}

	j.lg.Info("success to process job", zap.String("payload", payload))
	return nil
}
