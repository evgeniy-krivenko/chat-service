package managerv1

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	internalerrors "github.com/evgeniy-krivenko/chat-service/internal/errors"
	"github.com/evgeniy-krivenko/chat-service/internal/middlewares"
	sendmessage "github.com/evgeniy-krivenko/chat-service/internal/usecases/manager/send-message"
)

func (h Handlers) PostSendMessage(eCtx echo.Context, params PostSendMessageParams) error {
	ctx := eCtx.Request().Context()
	managerID := middlewares.MustUserID(eCtx)

	var req SendMessageRequest
	if err := eCtx.Bind(&req); err != nil {
		return internalerrors.NewServerError(http.StatusBadRequest, "bind request", err)
	}

	resp, err := h.sendMessageUseCase.Handle(ctx, sendmessage.Request{
		ID:          params.XRequestID,
		ManagerID:   managerID,
		ChatID:      req.ChatId,
		MessageBody: req.MessageBody,
	})
	if errors.Is(err, sendmessage.ErrInvalidRequest) {
		return internalerrors.NewServerError(http.StatusBadRequest, "bad request", err)
	}
	if errors.Is(err, sendmessage.ErrProblemNotFound) {
		return internalerrors.NewServerError(http.StatusBadRequest, "problem not found", err)
	}
	if err != nil {
		return internalerrors.NewServerError(http.StatusInternalServerError, "internal err", err)
	}

	response := SendMessageResponse{
		Data: &MessageWithoutBody{
			Id:        resp.MessageID,
			AuthorId:  managerID,
			CreatedAt: resp.CreatedAt,
		},
	}

	return eCtx.JSON(http.StatusOK, &response)
}
