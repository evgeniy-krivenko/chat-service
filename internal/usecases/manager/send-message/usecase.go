package sendmessage

import (
	"context"
	"errors"
	"fmt"
	"time"

	messagesrepo "github.com/evgeniy-krivenko/chat-service/internal/repositories/messages"
	problemsrepo "github.com/evgeniy-krivenko/chat-service/internal/repositories/problems"
	"github.com/evgeniy-krivenko/chat-service/internal/services/outbox/jobs/payload/simpleid"
	sendmanagermessagejob "github.com/evgeniy-krivenko/chat-service/internal/services/outbox/jobs/send-manager-message"
	"github.com/evgeniy-krivenko/chat-service/internal/types"
)

var (
	ErrInvalidRequest  = errors.New("invalid request")
	ErrProblemNotFound = errors.New("problem not found")
)

//go:generate mockgen -source=$GOFILE -destination=mocks/usecase_mock.gen.go -package=sendmessagemocks

type messagesRepository interface {
	CreateFullVisible(
		ctx context.Context,
		reqID types.RequestID,
		problemID types.ProblemID,
		chatID types.ChatID,
		authorID types.UserID,
		msgBody string,
	) (*messagesrepo.Message, error)
}

type outboxService interface {
	Put(ctx context.Context, name, payload string, availableAt time.Time) (types.JobID, error)
}

type problemsRepository interface {
	GetAssignedProblem(ctx context.Context, managerID types.UserID, chatID types.ChatID) (*problemsrepo.Problem, error)
}

type transactor interface {
	RunInTx(ctx context.Context, f func(context.Context) error) error
}

//go:generate options-gen -out-filename=usecase_options.gen.go -from-struct=Options
type Options struct {
	msgRepo      messagesRepository `option:"mandatory" validate:"required"`
	problemsRepo problemsRepository `option:"mandatory" validate:"required"`
	outboxSvc    outboxService      `option:"mandatory" validate:"required"`
	txtor        transactor         `option:"mandatory" validate:"required"`
}

type UseCase struct {
	Options
}

func New(opts Options) (UseCase, error) {
	return UseCase{Options: opts}, opts.Validate()
}

func (u UseCase) Handle(ctx context.Context, req Request) (Response, error) {
	if err := req.Validate(); err != nil {
		return Response{}, fmt.Errorf("%w: %v", ErrInvalidRequest, err)
	}

	problem, err := u.problemsRepo.GetAssignedProblem(ctx, req.ManagerID, req.ChatID)
	if err != nil {
		if errors.Is(err, problemsrepo.ErrProblemNotFound) {
			return Response{}, fmt.Errorf("%w: %v", ErrProblemNotFound, err)
		}
		return Response{}, fmt.Errorf("get assigned problem: %v", err)
	}

	var response Response

	if err := u.txtor.RunInTx(ctx, func(ctx context.Context) error {
		msg, err := u.msgRepo.CreateFullVisible(ctx, req.ID, problem.ID, req.ChatID, req.ManagerID, req.MessageBody)
		if err != nil {
			return fmt.Errorf("create full visible msg: %v", err)
		}

		_, err = u.outboxSvc.Put(ctx, sendmanagermessagejob.Name, simpleid.MustMarshal(msg.ID), time.Now())
		if err != nil {
			return fmt.Errorf("put job %v for msg %v: %v", sendmanagermessagejob.Name, msg.ID, err)
		}

		response.MessageID = msg.ID
		response.CreatedAt = msg.CreatedAt

		return nil
	}); err != nil {
		return Response{}, fmt.Errorf("save message: %v", err)
	}

	return response, nil
}
