package canreceiveproblems

import (
	"context"
	"errors"
	"fmt"

	"github.com/evgeniy-krivenko/chat-service/internal/types"
)

var ErrInvalidRequest = errors.New("invalid request")

//go:generate mockgen -source=$GOFILE -destination=mocks/usecase_mock.gen.go -package=canreceiveproblemsmocks

type managerLoadService interface {
	CanManagerTakeProblem(ctx context.Context, managerID types.UserID) (bool, error)
}

type managerPool interface {
	Contains(ctx context.Context, managerID types.UserID) (bool, error)
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
		return UseCase{}, fmt.Errorf("validate can receive problem usecase: %v", err)
	}
	return UseCase{Options: opts}, nil
}

func (u UseCase) Handle(ctx context.Context, req Request) (Response, error) {
	if err := req.Validate(); err != nil {
		return Response{}, fmt.Errorf("%w: %v", ErrInvalidRequest, err)
	}

	inPool, err := u.managerPool.Contains(ctx, req.ManagerID)
	if err != nil {
		return Response{}, fmt.Errorf("contains in manager pool with id %v: %v", req.ManagerID, err)
	}

	if inPool {
		return Response{Available: false, InPool: true}, nil
	}

	canTakeProblem, err := u.managerLoadSrv.CanManagerTakeProblem(ctx, req.ManagerID)
	if err != nil {
		return Response{}, fmt.Errorf("availability in manager load service with id %v: %v", req.ManagerID, err)
	}

	return Response{Available: canTakeProblem, InPool: false}, nil
}
