package managerassignedtoproblemjob

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	chatsrepo "github.com/evgeniy-krivenko/chat-service/internal/repositories/chats"
	messagesrepo "github.com/evgeniy-krivenko/chat-service/internal/repositories/messages"
	profilesrepo "github.com/evgeniy-krivenko/chat-service/internal/repositories/profiles"
	eventstream "github.com/evgeniy-krivenko/chat-service/internal/services/event-stream"
	"github.com/evgeniy-krivenko/chat-service/internal/services/outbox"
	"github.com/evgeniy-krivenko/chat-service/internal/types"
)

const Name = "manager-assigned-to-problem"

//go:generate mockgen -source=$GOFILE -destination=mocks/job_mock.gen.go -package=managerassignedtoproblemjobmocks

type messageRepo interface {
	GetMessageByID(ctx context.Context, msgID types.MessageID) (*messagesrepo.Message, error)
}

type chatRepo interface {
	GetChatByID(ctx context.Context, chatID types.ChatID) (*chatsrepo.Chat, error)
}

type profilesRepo interface {
	GetProfileByID(ctx context.Context, id types.UserID) (profile *profilesrepo.Profile, err error)
}

type managerLoadService interface {
	CanManagerTakeProblem(ctx context.Context, managerID types.UserID) (bool, error)
}

type eventStream interface {
	Publish(ctx context.Context, userID types.UserID, event eventstream.Event) error
}

//go:generate options-gen -out-filename=job_options.gen.go -from-struct=Options
type Options struct {
	chatRepo     chatRepo           `option:"mandatory" validate:"required"`
	msgRepo      messageRepo        `option:"mandatory" validate:"required"`
	profilesRepo profilesRepo       `option:"mandatory" validate:"required"`
	mngrLoadSvc  managerLoadService `option:"mandatory" validate:"required"`
	eventStream  eventStream        `option:"mandatory" validate:"required"`
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
	p, err := unmarshalPayload(payload)
	if err != nil {
		j.lg.Error("unmarshal payload", zap.Error(err))
		return fmt.Errorf("unmarshal payload: %v", err)
	}

	msg, err := j.msgRepo.GetMessageByID(ctx, p.MessageID)
	if err != nil {
		j.lg.Error("get msg from repo", zap.Error(err))

		return fmt.Errorf("get msg while hadle job %v: %v", Name, err)
	}

	chat, err := j.chatRepo.GetChatByID(ctx, msg.ChatID)
	if err != nil {
		j.lg.Warn("get chat", zap.Stringer("chat_id", msg.ChatID), zap.Error(err))

		return fmt.Errorf("get chat %v: %v", msg.ChatID, err)
	}

	profile, err := j.profilesRepo.GetProfileByID(ctx, chat.ClientID)
	if err != nil {
		j.lg.Warn("not find profile", zap.Stringer("user_id", chat.ClientID), zap.Error(err))

		return fmt.Errorf("get profile: %v", err)
	}

	canTakeMoreProblems, err := j.mngrLoadSvc.CanManagerTakeProblem(ctx, p.ManagerID)
	if err != nil {
		j.lg.Warn("can take problem", zap.Stringer("manager_id", p.ManagerID), zap.Error(err))

		return fmt.Errorf("can manager %v take problem: %v", p.ManagerID, err)
	}

	if err := j.eventStream.Publish(ctx, chat.ClientID, eventstream.NewNewMessageEvent(
		types.NewEventID(),
		p.RequestID,
		msg.ChatID,
		msg.ID,
		types.UserIDNil,
		msg.CreatedAt,
		msg.Body,
		"",
		msg.IsService,
	)); err != nil {
		j.lg.Warn(
			"send new message event to client",
			zap.Stringer("client_id", chat.ClientID),
			zap.Error(err),
		)

		return fmt.Errorf("send new message event to client %v: %v", chat.ClientID, err)
	}

	if err := j.eventStream.Publish(ctx, p.ManagerID, eventstream.NewNewChatEvent(
		types.NewEventID(),
		p.RequestID,
		chat.ID,
		chat.ClientID,
		profile.FirstName,
		profile.LastName,
		canTakeMoreProblems,
	)); err != nil {
		j.lg.Warn("publish new chat event to manager", zap.Stringer("manager_id", p.ManagerID), zap.Error(err))

		return fmt.Errorf("publish new chat event to manager %v: %v", p.ManagerID, err)
	}

	return nil
}
