package getchathistory

import (
	"context"
	"errors"
	"fmt"

	"github.com/evgeniy-krivenko/chat-service/internal/cursor"
	messagesrepo "github.com/evgeniy-krivenko/chat-service/internal/repositories/messages"
	problemsrepo "github.com/evgeniy-krivenko/chat-service/internal/repositories/problems"
	"github.com/evgeniy-krivenko/chat-service/internal/types"
	"github.com/evgeniy-krivenko/chat-service/pkg/utils"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/usecase_mock.gen.go -package=getchathistorymocks

var (
	ErrInvalidRequest = errors.New("invalid request")
	ErrInvalidCursor  = errors.New("invalid cursor")
)

type messagesRepository interface {
	GetProblemMessages(
		ctx context.Context,
		problemID types.ProblemID,
		pageSize int,
		cursor *messagesrepo.Cursor,
	) ([]messagesrepo.Message, *messagesrepo.Cursor, error)
}

type problemsRepository interface {
	GetOpenProblemForChat(
		ctx context.Context,
		chatID types.ChatID,
		managerID types.UserID,
	) (*problemsrepo.Problem, error)
}

//go:generate options-gen -out-filename=usecase_options.gen.go -from-struct=Options
type Options struct {
	msgRepo      messagesRepository `option:"mandatory" validate:"required"`
	problemsRepo problemsRepository `option:"mandatory" validate:"required"`
}

type UseCase struct {
	Options
}

func New(opts Options) (UseCase, error) {
	if err := opts.Validate(); err != nil {
		return UseCase{}, fmt.Errorf("validate gethistory usecase: %v", err)
	}
	return UseCase{opts}, nil
}

func (u UseCase) Handle(ctx context.Context, req Request) (Response, error) {
	var response Response

	if err := req.Validate(); err != nil {
		return Response{}, fmt.Errorf("%w: %v", ErrInvalidRequest, err)
	}

	var decodedCursor *messagesrepo.Cursor
	if req.Cursor != "" {
		if err := cursor.Decode(req.Cursor, &decodedCursor); err != nil {
			return Response{}, fmt.Errorf("%w: %v", ErrInvalidCursor, err)
		}
	}

	problem, err := u.problemsRepo.GetOpenProblemForChat(ctx, req.ChatID, req.ManagerID)
	if err != nil {
		return Response{}, fmt.Errorf(
			"get problem for manager %v and chat %v: %v",
			req.ManagerID,
			req.ChatID,
			err,
		)
	}

	messages, newCursor, err := u.msgRepo.GetProblemMessages(ctx, problem.ID, req.PageSize, decodedCursor)
	if err != nil {
		if errors.Is(err, messagesrepo.ErrInvalidCursor) {
			err = fmt.Errorf("%w: %v", ErrInvalidCursor, err)
		}

		return Response{}, err
	}

	response.Messages = utils.Apply(messages, adaptMessage)

	if newCursor != nil {
		response.NextCursor, err = cursor.Encode(newCursor)
		if err != nil {
			return Response{}, err
		}
	}

	return response, nil
}
