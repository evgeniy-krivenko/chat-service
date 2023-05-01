package clientv1

import (
	"context"
	"fmt"

	gethistory "github.com/evgeniy-krivenko/chat-service/internal/usecases/client/get-history"
	sendmessage "github.com/evgeniy-krivenko/chat-service/internal/usecases/client/send-message"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/handlers_mocks.gen.go -package=clientv1mocks

type getHistoryUseCase interface {
	Handle(ctx context.Context, req gethistory.Request) (gethistory.Response, error)
}

type sendMessageUseCase interface {
	Handle(ctx context.Context, req sendmessage.Request) (sendmessage.Response, error)
}

//go:generate options-gen -out-filename=handler_options.gen.go -from-struct=Options
type Options struct {
	getHistory     getHistoryUseCase  `option:"mandatory" validate:"required"`
	sendMsgUseCase sendMessageUseCase `option:"mandatory" validate:"required"`
}

type Handlers struct {
	Options
}

func NewHandlers(opts Options) (Handlers, error) {
	if err := opts.Validate(); err != nil {
		return Handlers{}, fmt.Errorf("validate handlers options: %v", err)
	}

	return Handlers{Options: opts}, nil
}
