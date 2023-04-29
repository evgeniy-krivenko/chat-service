package sendmessage

import (
	"context"
	"errors"
	"time"

	messagesrepo "github.com/evgeniy-krivenko/chat-service/internal/repositories/messages"
	sendclientmessagejob "github.com/evgeniy-krivenko/chat-service/internal/services/outbox/jobs/send-client-message"
	"github.com/evgeniy-krivenko/chat-service/internal/types"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/usecase_mock.gen.go -package=sendmessagemocks

var (
	ErrInvalidRequest    = errors.New("invalid request")
	ErrChatNotCreated    = errors.New("chat not created")
	ErrProblemNotCreated = errors.New("problem not created")
)

type chatsRepository interface {
	CreateIfNotExists(ctx context.Context, userID types.UserID) (types.ChatID, error)
}

type outboxService interface {
	Put(ctx context.Context, name, payload string, availableAt time.Time) (types.JobID, error)
}

type messagesRepository interface {
	GetMessageByRequestID(ctx context.Context, reqID types.RequestID) (*messagesrepo.Message, error)
	CreateClientVisible(
		ctx context.Context,
		reqID types.RequestID,
		problemID types.ProblemID,
		chatID types.ChatID,
		authorID types.UserID,
		msgBody string,
	) (*messagesrepo.Message, error)
}

type problemsRepository interface {
	CreateIfNotExists(ctx context.Context, chatID types.ChatID) (types.ProblemID, error)
}

type transactor interface {
	RunInTx(ctx context.Context, f func(context.Context) error) error
}

//go:generate options-gen -out-filename=usecase_options.gen.go -from-struct=Options
type Options struct {
	chatRepo    chatsRepository    `option:"mandatory" validate:"required"`
	msgRepo     messagesRepository `option:"mandatory" validate:"required"`
	outboxSrv   outboxService      `option:"mandatory" validate:"required"`
	problemRepo problemsRepository `option:"mandatory" validate:"required"`
	txtor       transactor         `option:"mandatory" validate:"required"`
}

type UseCase struct {
	Options
}

func New(opts Options) (UseCase, error) {
	return UseCase{Options: opts}, opts.Validate()
}

func (u UseCase) Handle(ctx context.Context, req Request) (Response, error) {
	if err := req.Validate(); err != nil {
		return Response{}, ErrInvalidRequest
	}

	var msg *messagesrepo.Message

	err := u.txtor.RunInTx(ctx, func(ctx context.Context) error {
		m, err := u.msgRepo.GetMessageByRequestID(ctx, req.ID)
		if errors.Is(err, messagesrepo.ErrMsgNotFound) {
			chatID, err := u.chatRepo.CreateIfNotExists(ctx, req.ClientID)
			if err != nil {
				return ErrChatNotCreated
			}

			problemID, err := u.problemRepo.CreateIfNotExists(ctx, chatID)
			if err != nil {
				return ErrProblemNotCreated
			}

			newMsg, err := u.msgRepo.CreateClientVisible(ctx, req.ID, problemID, chatID, req.ClientID, req.MessageBody)
			if err != nil {
				return err
			}

			msg = newMsg

			return u.putToOutbox(ctx, newMsg.ID)
		}
		if err != nil {
			return err
		}

		msg = m
		return nil
	})
	if err != nil {
		return Response{}, err
	}

	return Response{
		MessageID: msg.ID,
		AuthorID:  msg.AuthorID,
		CreatedAt: msg.CreatedAt,
	}, nil
}

func (u UseCase) putToOutbox(ctx context.Context, msgID types.MessageID) error {
	outboxPayload, err := sendclientmessagejob.MarshalPayload(msgID)
	if err != nil {
		return err
	}

	_, err = u.outboxSrv.Put(ctx, sendclientmessagejob.Name, outboxPayload, time.Now())
	if err != nil {
		return err
	}

	return nil
}
