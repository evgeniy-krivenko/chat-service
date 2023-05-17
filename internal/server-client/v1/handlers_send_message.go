package clientv1

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	internalerrors "github.com/evgeniy-krivenko/chat-service/internal/errors"
	"github.com/evgeniy-krivenko/chat-service/internal/middlewares"
	sendmessage "github.com/evgeniy-krivenko/chat-service/internal/usecases/client/send-message"
	"github.com/evgeniy-krivenko/chat-service/pkg/pointer"
)

func (h Handlers) PostSendMessage(eCtx echo.Context, params PostSendMessageParams) error {
	ctx := eCtx.Request().Context()
	clientID := middlewares.MustUserID(eCtx)

	var req SendMessageRequest
	if err := eCtx.Bind(&req); err != nil {
		return internalerrors.NewServerError(http.StatusBadRequest, "bind request", err)
	}

	message, err := h.sendMsgUseCase.Handle(ctx, sendmessage.Request{
		ID:          params.XRequestID,
		ClientID:    clientID,
		MessageBody: req.MessageBody,
	})
	if err != nil {
		if errors.Is(err, sendmessage.ErrInvalidRequest) {
			return internalerrors.NewServerError(http.StatusBadRequest, "invalid request", err)
		}
		if errors.Is(err, sendmessage.ErrChatNotCreated) {
			return internalerrors.NewServerError(ErrorCodeCreateChatError, "create chat", err)
		}
		if errors.Is(err, sendmessage.ErrProblemNotCreated) {
			return internalerrors.NewServerError(ErrorCodeCreateProblemError, "create problem", err)
		}

		return internalerrors.NewServerError(http.StatusInternalServerError, "unknown error while send message", err)
	}

	return eCtx.JSON(http.StatusOK, &SendMessageResponse{
		Data: &MessageHeader{
			AuthorId:  pointer.PtrWithZeroAsNil(message.AuthorID),
			CreatedAt: message.CreatedAt,
			Id:        message.MessageID,
		},
	})
}
