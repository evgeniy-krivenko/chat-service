package closechat

import (
	"context"
	"errors"
	"fmt"
	"time"

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
	Resolve(ctx context.Context, managerID types.UserID, chatID types.ChatID) error
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

	payload, err := closechatjob.MarshalPayload(req.ID, req.ManagerID, req.ChatID)
	if err != nil {
		return fmt.Errorf("marshal payload: %v", err)
	}

	return u.txtor.RunInTx(ctx, func(ctx context.Context) error {
		if err := u.problemsRepo.Resolve(ctx, req.ManagerID, req.ChatID); err != nil {
			if errors.Is(err, problemsrepo.ErrProblemNotFound) {
				return ErrProblemNotFound
			}
			return fmt.Errorf("resolve problem: %v", err)
		}

		_, err = u.outboxSvc.Put(ctx, closechatjob.Name, payload, time.Now())
		if err != nil {
			return fmt.Errorf("put job: %v", err)
		}

		return nil
	})
}
