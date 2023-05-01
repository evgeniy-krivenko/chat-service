package gethistory

import (
	"context"
	"errors"
	"fmt"

	"github.com/evgeniy-krivenko/chat-service/internal/cursor"
	messagesrepo "github.com/evgeniy-krivenko/chat-service/internal/repositories/messages"
	"github.com/evgeniy-krivenko/chat-service/internal/types"
	"github.com/evgeniy-krivenko/chat-service/pkg/utils"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/usecase_mock.gen.go -package=gethistorymocks

var (
	ErrInvalidRequest = errors.New("invalid request")
	ErrInvalidCursor  = errors.New("invalid cursor")
)

type messagesRepository interface {
	GetClientChatMessages(
		ctx context.Context,
		clientID types.UserID,
		pageSize int,
		cursor *messagesrepo.Cursor,
	) ([]messagesrepo.Message, *messagesrepo.Cursor, error)
}

//go:generate options-gen -out-filename=usecase_options.gen.go -from-struct=Options
type Options struct {
	msgRepo messagesRepository `option:"mandatory" validate:"required"`
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
		return Response{}, ErrInvalidRequest
	}

	var decodedCursor *messagesrepo.Cursor
	if req.Cursor != "" {
		if err := cursor.Decode(req.Cursor, &decodedCursor); err != nil {
			return Response{}, fmt.Errorf("%w: %v", ErrInvalidCursor, err)
		}
	}

	messages, newCursor, err := u.msgRepo.GetClientChatMessages(ctx, req.ClientID, req.PageSize, decodedCursor)
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
