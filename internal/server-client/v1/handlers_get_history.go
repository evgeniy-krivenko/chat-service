package clientv1

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	internalerrors "github.com/evgeniy-krivenko/chat-service/internal/errors"
	"github.com/evgeniy-krivenko/chat-service/internal/middlewares"
	gethistory "github.com/evgeniy-krivenko/chat-service/internal/usecases/client/get-history"
	"github.com/evgeniy-krivenko/chat-service/pkg/pointer"
	"github.com/evgeniy-krivenko/chat-service/pkg/utils"
)

func (h Handlers) PostGetHistory(
	eCtx echo.Context,
	params PostGetHistoryParams,
) error {
	ctx := eCtx.Request().Context()
	clientID := middlewares.MustUserID(eCtx)

	var req GetHistoryRequest
	err := eCtx.Bind(&req)
	if err != nil {
		return internalerrors.NewServerError(http.StatusBadRequest, "bind request", err)
	}

	useCaseResponse, err := h.getHistory.Handle(ctx, gethistory.Request{
		ID:       params.XRequestID,
		ClientID: clientID,
		PageSize: pointer.Indirect(req.PageSize),
		Cursor:   pointer.Indirect(req.Cursor),
	})
	if err != nil {
		if errors.Is(err, gethistory.ErrInvalidRequest) {
			return internalerrors.NewServerError(http.StatusBadRequest, "invalid request for get history", err)
		}
		if errors.Is(err, gethistory.ErrInvalidCursor) {
			return internalerrors.NewServerError(http.StatusBadRequest, "invalid cursor for get history", err)
		}

		return internalerrors.NewServerError(http.StatusInternalServerError, "unknown error while get history", err)
	}

	page := &MessagesPage{
		Messages: utils.Apply(useCaseResponse.Messages, adaptGetHistoryMsg),
		Next:     useCaseResponse.NextCursor,
	}

	return eCtx.JSON(http.StatusOK, &GetHistoryResponse{Data: page})
}

func adaptGetHistoryMsg(msg gethistory.Message) Message {
	return Message{
		Id:         msg.ID,
		AuthorId:   pointer.PtrWithZeroAsNil(msg.AuthorID),
		Body:       msg.Body,
		CreatedAt:  msg.CreatedAt,
		IsBlocked:  msg.IsBlocked,
		IsReceived: msg.IsReceived,
		IsService:  msg.IsService,
		AuthorName: pointer.PtrWithZeroAsNil(msg.AuthorName),
	}
}
