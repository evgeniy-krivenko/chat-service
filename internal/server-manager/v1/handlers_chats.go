package managerv1

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	internalerrors "github.com/evgeniy-krivenko/chat-service/internal/errors"
	"github.com/evgeniy-krivenko/chat-service/internal/middlewares"
	getchathistory "github.com/evgeniy-krivenko/chat-service/internal/usecases/manager/get-chat-history"
	getchats "github.com/evgeniy-krivenko/chat-service/internal/usecases/manager/get-chats"
	"github.com/evgeniy-krivenko/chat-service/pkg/pointer"
	"github.com/evgeniy-krivenko/chat-service/pkg/utils"
)

func (h Handlers) PostGetChats(eCtx echo.Context, params PostGetChatsParams) error {
	ctx := eCtx.Request().Context()
	managerID := middlewares.MustUserID(eCtx)

	response, err := h.getChatsUseCase.Handle(ctx, getchats.Request{
		ID:        params.XRequestID,
		ManagerID: managerID,
	})
	if errors.Is(err, getchats.ErrInvalidRequest) {
		return internalerrors.NewServerError(http.StatusBadRequest, "invalid request", err)
	}
	if err != nil {
		return internalerrors.NewServerError(http.StatusInternalServerError, "internal error", err)
	}

	adaptChats := utils.Apply(response.Chats, adaptGetChats)

	return eCtx.JSON(http.StatusOK, &GetChatsResponse{Data: &ChatList{
		adaptChats,
	}})
}

func (h Handlers) PostGetChatHistory(eCtx echo.Context, params PostGetChatHistoryParams) error {
	ctx := eCtx.Request().Context()
	managerID := middlewares.MustUserID(eCtx)

	var req GetChatHistoryRequest
	if err := eCtx.Bind(&req); err != nil {
		return internalerrors.NewServerError(http.StatusBadRequest, "bind request", err)
	}

	useCaseResponse, err := h.getChatHistoryUseCase.Handle(ctx, getchathistory.Request{
		ID:        params.XRequestID,
		ManagerID: managerID,
		ChatID:    req.ChatId,
		PageSize:  pointer.Indirect(req.PageSize),
		Cursor:    pointer.Indirect(req.Cursor),
	})
	if err != nil {
		if errors.Is(err, getchathistory.ErrInvalidRequest) {
			return internalerrors.NewServerError(http.StatusBadRequest, "invalid request for get history", err)
		}
		if errors.Is(err, getchathistory.ErrInvalidCursor) {
			return internalerrors.NewServerError(http.StatusBadRequest, "invalid cursor for get history", err)
		}

		return internalerrors.NewServerError(http.StatusInternalServerError, "unknown error while get history", err)
	}

	page := &MessagesPage{
		Messages: utils.Apply(useCaseResponse.Messages, adaptGetChatHistoryMsg),
		Next:     useCaseResponse.NextCursor,
	}

	return eCtx.JSON(http.StatusOK, &GetChatHistoryResponse{Data: page})
}

func adaptGetChats(c getchats.Chat) Chat {
	return Chat{
		ChatId:    c.ID,
		ClientId:  c.ClientID,
		FirstName: pointer.PtrWithZeroAsNil(c.FirstName),
		LastName:  pointer.PtrWithZeroAsNil(c.LastName),
	}
}

func adaptGetChatHistoryMsg(m getchathistory.Message) Message {
	return Message{
		Id:        m.ID,
		Body:      m.Body,
		CreatedAt: m.CreatedAt,
		AuthorId:  m.AuthorID,
	}
}
