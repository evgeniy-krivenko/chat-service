package managerv1

import (
	"context"
	"fmt"

	canreceiveproblems "github.com/evgeniy-krivenko/chat-service/internal/usecases/manager/can-receive-problems"
	closechat "github.com/evgeniy-krivenko/chat-service/internal/usecases/manager/close-chat"
	freehands "github.com/evgeniy-krivenko/chat-service/internal/usecases/manager/free-hands"
	getchathistory "github.com/evgeniy-krivenko/chat-service/internal/usecases/manager/get-chat-history"
	getchats "github.com/evgeniy-krivenko/chat-service/internal/usecases/manager/get-chats"
	sendmessage "github.com/evgeniy-krivenko/chat-service/internal/usecases/manager/send-message"
)

var _ ServerInterface = (*Handlers)(nil)

//go:generate mockgen -source=$GOFILE -destination=mocks/handlers_mocks.gen.go -package=managerv1mocks

type canReceiveProblemsUseCase interface {
	Handle(ctx context.Context, req canreceiveproblems.Request) (canreceiveproblems.Response, error)
}

type freeHandsUseCase interface {
	Handle(ctx context.Context, req freehands.Request) error
}

type getChatsUseCase interface {
	Handle(ctx context.Context, req getchats.Request) (getchats.Response, error)
}

type getChatHistoryUseCase interface {
	Handle(ctx context.Context, req getchathistory.Request) (getchathistory.Response, error)
}

type sendMessageUseCase interface {
	Handle(ctx context.Context, req sendmessage.Request) (sendmessage.Response, error)
}

type closeChatUseCase interface {
	Handle(ctx context.Context, req closechat.Request) error
}

//go:generate options-gen -out-filename=handler_options.gen.go -from-struct=Options
type Options struct {
	canReceiveProblemUseCase canReceiveProblemsUseCase `option:"mandatory" validate:"required"`
	freeHandsUseCase         freeHandsUseCase          `option:"mandatory" validate:"required"`
	getChatsUseCase          getChatsUseCase           `option:"mandatory" validate:"required"`
	getChatHistoryUseCase    getChatHistoryUseCase     `option:"mandatory" validate:"required"`
	sendMessageUseCase       sendMessageUseCase        `option:"mandatory" validate:"required"`
	closeChatUseCase         closeChatUseCase          `option:"mandatory" validate:"required"`
}

type Handlers struct {
	Options
}

func NewHandlers(opts Options) (Handlers, error) {
	if err := opts.Validate(); err != nil {
		return Handlers{}, fmt.Errorf("validate manager handlers: %v", err)
	}
	return Handlers{Options: opts}, nil
}
