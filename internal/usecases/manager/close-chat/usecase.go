package closechat

import (
	"context"
	"errors"
	"fmt"
	"time"

	messagesrepo "github.com/evgeniy-krivenko/chat-service/internal/repositories/messages"
	problemsrepo "github.com/evgeniy-krivenko/chat-service/internal/repositories/problems"
	closechatjob "github.com/evgeniy-krivenko/chat-service/internal/services/outbox/jobs/close-chat"
	"github.com/evgeniy-krivenko/chat-service/internal/types"
)

var (
	ErrInvalidRequest  = errors.New("invalid request")
	ErrProblemNotFound = errors.New("problem not found")
)

//go:generate mockgen -source=$GOFILE -destination=mocks/usecase_mock.gen.go -package=closechatmocks

type problemsRepository interface {
	GetAssignedProblem(
		ctx context.Context,
		managerID types.UserID,
		chatID types.ChatID,
	) (*problemsrepo.Problem, error)
	Resolve(
		ctx context.Context,
		reqID types.RequestID,
		managerID types.UserID,
		chatID types.ChatID,
	) error
}

type messageRepository interface {
	CreateClientService(
		ctx context.Context,
		problemID types.ProblemID,
		chatID types.ChatID,
		msgBody string,
	) (*messagesrepo.Message, error)
}

type outboxService interface {
	Put(ctx context.Context, name, payload string, availableAt time.Time) (types.JobID, error)
}

type transactor interface {
	RunInTx(ctx context.Context, f func(context.Context) error) error
}

//go:generate options-gen -out-filename=usecase_options.gen.go -from-struct=Options
type Options struct {
	problemsRepo problemsRepository `option:"mandatory" validate:"required"`
	msgRepo      messageRepository  `option:"mandatory" validate:"required"`
	outboxSvc    outboxService      `option:"mandatory" validate:"required"`
	txtor        transactor         `option:"mandatory" validate:"required"`
}

type UseCase struct {
	Options
}

func New(opts Options) (UseCase, error) {
	return UseCase{Options: opts}, opts.Validate()
}

func (u UseCase) Handle(ctx context.Context, req Request) error {
	if err := req.Validate(); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidRequest, err)
	}

	return u.txtor.RunInTx(ctx, func(ctx context.Context) error {
		problem, err := u.problemsRepo.GetAssignedProblem(ctx, req.ManagerID, req.ChatID)
		if err != nil {
			if errors.Is(err, problemsrepo.ErrProblemNotFound) {
				return ErrProblemNotFound
			}
			return fmt.Errorf("get assigned problem: %v", err)
		}

		if err := u.problemsRepo.Resolve(ctx, req.ID, req.ManagerID, req.ChatID); err != nil {
			return fmt.Errorf("resolve problem: %v", err)
		}

		msg, err := u.msgRepo.CreateClientService(
			ctx,
			problem.ID,
			req.ChatID,
			"Your question has been marked as resolved.\nThank you for being with us!",
		)
		if err != nil {
			return fmt.Errorf("create client service problem")
		}

		payload, err := closechatjob.MarshalPayload(req.ID, req.ManagerID, req.ChatID, msg.ID)
		if err != nil {
			return fmt.Errorf("marshal payload: %v", err)
		}

		_, err = u.outboxSvc.Put(ctx, closechatjob.Name, payload, time.Now())
		if err != nil {
			return fmt.Errorf("put job: %v", err)
		}

		return nil
	})
}
