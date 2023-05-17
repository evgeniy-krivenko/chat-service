package freehands

import (
	"context"
	"errors"
	"fmt"

	"github.com/evgeniy-krivenko/chat-service/internal/types"
)

var ErrManagerOverloaded = errors.New("manager overloaded")

//go:generate mockgen -source=$GOFILE -destination=mocks/usecase_mock.gen.go -package=freehandsmocks

type managerLoadService interface {
	CanManagerTakeProblem(ctx context.Context, managerID types.UserID) (bool, error)
}

type managerPool interface {
	Put(ctx context.Context, managerID types.UserID) error
}

//go:generate options-gen -out-filename=usecase_options.gen.go -from-struct=Options
type Options struct {
	managerLoadSrv managerLoadService `option:"mandatory" validate:"required"`
	managerPool    managerPool        `option:"mandatory" validate:"required"`
}

type UseCase struct {
	Options
}

func New(opts Options) (UseCase, error) {
	if err := opts.Validate(); err != nil {
		return UseCase{}, fmt.Errorf("validate free hands usecase: %v", err)
	}
	return UseCase{Options: opts}, nil
}

func (u UseCase) Handle(ctx context.Context, req Request) error {
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validate free hands usecase options: %v", err)
	}

	canTakeProblem, err := u.managerLoadSrv.CanManagerTakeProblem(ctx, req.ManagerID)
	if err != nil {
		return fmt.Errorf("free hands availability in manager load svr with id %v: %v", req.ManagerID, err)
	}

	if !canTakeProblem {
		return ErrManagerOverloaded
	}

	err = u.managerPool.Put(ctx, req.ManagerID)
	if err != nil {
		return fmt.Errorf("put manager in pool with id %v: %v", req.ManagerID, err)
	}

	return nil
}
