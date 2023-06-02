package closechatjob

import (
	"context"
	"errors"
	"fmt"
	messagesrepo "github.com/evgeniy-krivenko/chat-service/internal/repositories/messages"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	chatsrepo "github.com/evgeniy-krivenko/chat-service/internal/repositories/chats"
	eventstream "github.com/evgeniy-krivenko/chat-service/internal/services/event-stream"
	"github.com/evgeniy-krivenko/chat-service/internal/services/outbox"
	"github.com/evgeniy-krivenko/chat-service/internal/types"
)

const Name = "close-chat"

const CloseMsgBody = "Your question has been marked as resolved.\nThank you for being with us!"

//go:generate mockgen -source=$GOFILE -destination=mocks/job_mock.gen.go -package=closechatjobmocks

type managerLoadService interface {
	CanManagerTakeProblem(ctx context.Context, managerID types.UserID) (bool, error)
}

type eventStream interface {
	Publish(ctx context.Context, userID types.UserID, event eventstream.Event) error
}

type chatsRepository interface {
	GetChatByID(ctx context.Context, chatID types.ChatID) (*chatsrepo.Chat, error)
}

type messageRepository interface {
	GetMessageByID(ctx context.Context, msgID types.MessageID) (*messagesrepo.Message, error)
}

//go:generate options-gen -out-filename=job_options.gen.go -from-struct=Options
type Options struct {
	mngrLoadSvc managerLoadService `option:"mandatory" validate:"required"`
	eventStream eventStream        `option:"mandatory" validate:"required"`
	chatsRepo   chatsRepository    `option:"mandatory" validate:"required"`
	msgRepo     messageRepository  `option:"mandatory" validate:"required"`
}

type Job struct {
	outbox.DefaultJob
	Options
	lg *zap.Logger
}

func New(opts Options) (*Job, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("validate close chat job: %v", err)
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

	p, err := unmarshalPayload(payload)
	if err != nil {
		return fmt.Errorf("unmarshal payload: %v", err)
	}

	canManagerTakeProblem, err := j.mngrLoadSvc.CanManagerTakeProblem(ctx, p.ManagerID)
	if err != nil {
		return fmt.Errorf("can manager take problem: %v", err)
	}

	msg, err := j.msgRepo.GetMessageByID(ctx, p.ClientMsgID)
	if err != nil {
		return fmt.Errorf("get msg: %v", err)
	}

	chat, err := j.chatsRepo.GetChatByID(ctx, msg.ChatID)
	if err != nil {
		return fmt.Errorf("get chat: %v", err)
	}

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		if err := j.eventStream.Publish(ctx, chat.ClientID, eventstream.NewNewMessageEvent(
			types.NewEventID(),
			p.RequestID,
			msg.ChatID,
			msg.ID,
			msg.AuthorID,
			msg.CreatedAt,
			msg.Body,
			msg.IsService,
		)); err != nil {
			return fmt.Errorf("publish NewMesaggeEvent to client stream: %v", err)
		}

		return nil
	})

	eg.Go(func() error {
		if err := j.eventStream.Publish(ctx, p.ManagerID, eventstream.NewChatClosedEvent(
			types.NewEventID(),
			p.RequestID,
			p.ChatID,
			canManagerTakeProblem,
		)); err != nil {
			return fmt.Errorf("publish ChatClosedEvent to manager stream: %v", err)
		}

		return nil
	})

	if err := eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		j.lg.Warn("error to handle job", zap.String("payload", payload), zap.Error(err))
		return err
	}

	j.lg.Info("success to handle job", zap.String("payload", payload))
	return nil
}
