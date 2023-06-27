package getchats

import (
	"context"
	"errors"
	"fmt"

	chatsrepo "github.com/evgeniy-krivenko/chat-service/internal/repositories/chats"
	"github.com/evgeniy-krivenko/chat-service/internal/types"
	"github.com/evgeniy-krivenko/chat-service/pkg/utils"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/usecase_mock.gen.go -package=getchatsmocks

var ErrInvalidRequest = errors.New("invalid request")

type chatsRepository interface {
	GetManagerChatsWithProblems(ctx context.Context, managerID types.UserID) ([]chatsrepo.Chat, error)
}

//go:generate options-gen -out-filename=usecase_options.gen.go -from-struct=Options
type Options struct {
	chatsRepo chatsRepository `option:"mandatory" validate:"required"`
}

type UseCase struct {
	Options
}

func New(opts Options) (UseCase, error) {
	if err := opts.Validate(); err != nil {
		return UseCase{}, fmt.Errorf("vaildate get chats use case options: %v", err)
	}
	return UseCase{Options: opts}, nil
}

func (u UseCase) Handle(ctx context.Context, req Request) (Response, error) {
	if err := req.Validate(); err != nil {
		return Response{}, fmt.Errorf("%w: validate req for get chats usecase: %v", ErrInvalidRequest, err)
	}

	chats, err := u.chatsRepo.GetManagerChatsWithProblems(ctx, req.ManagerID)
	if err != nil {
		return Response{}, fmt.Errorf("get manager chats err: %w", err)
	}

	response := Response{
		Chats: utils.Apply(chats, adaptChat),
	}
	return response, nil
}

func adaptChat(chat chatsrepo.Chat) Chat {
	return Chat{
		ID:        chat.ID,
		ClientID:  chat.ClientID,
		FirstName: chat.FirstName,
		LastName:  chat.LastName,
	}
}
