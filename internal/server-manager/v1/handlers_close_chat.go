package managerv1

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	internalerrors "github.com/evgeniy-krivenko/chat-service/internal/errors"
	"github.com/evgeniy-krivenko/chat-service/internal/middlewares"
	closechat "github.com/evgeniy-krivenko/chat-service/internal/usecases/manager/close-chat"
)

func (h Handlers) PostCloseChat(eCtx echo.Context, params PostCloseChatParams) error {
	ctx := eCtx.Request().Context()
	managerID := middlewares.MustUserID(eCtx)

	var req ChatId
	if err := eCtx.Bind(&req); err != nil {
		return internalerrors.NewServerError(http.StatusBadRequest, "bind request", err)
	}

	if err := h.closeChatUseCase.Handle(ctx, closechat.Request{
		ID:        params.XRequestID,
		ManagerID: managerID,
		ChatID:    req.ChatId,
	}); err != nil {
		if errors.Is(err, closechat.ErrInvalidRequest) {
			return internalerrors.NewServerError(http.StatusBadRequest, "invalid request", err)
		}
		if errors.Is(err, closechat.ErrProblemNotFound) {
			return internalerrors.NewServerError(ErrorCodeNoFoundProblem, "no found open problem for chat", err)
		}
		return internalerrors.NewServerError(http.StatusInternalServerError, "internal error", err)
	}

	return eCtx.JSON(http.StatusOK, &CloseChatResponse{})
}
