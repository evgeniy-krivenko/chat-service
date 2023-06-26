package managerv1_test

import (
	"errors"
	"fmt"
	"net/http"

	internalerrors "github.com/evgeniy-krivenko/chat-service/internal/errors"
	managerv1 "github.com/evgeniy-krivenko/chat-service/internal/server-manager/v1"
	"github.com/evgeniy-krivenko/chat-service/internal/types"
	closechat "github.com/evgeniy-krivenko/chat-service/internal/usecases/manager/close-chat"
)

func (s *HandlersSuite) TestCloseChat_BindRequestError() {
	// Arrange.
	reqID := types.NewRequestID()
	resp, eCtx := s.newEchoCtx(reqID, "/v1/closeChat", `{"chatId":"`)

	// Action.
	err := s.handlers.PostCloseChat(eCtx, managerv1.PostCloseChatParams{XRequestID: reqID})

	// Assert.
	s.Require().Error(err)
	s.Equal(http.StatusBadRequest, internalerrors.GetServerErrorCode(err))
	s.Empty(resp.Body)
}

func (s *HandlersSuite) TestCloseChat_UseCase_InvalidRequestError() {
	// Arrange.
	reqID := types.NewRequestID()
	chatID := types.NewChatID()

	resp, eCtx := s.newEchoCtx(reqID, "/v1/closeChat", fmt.Sprintf(`{"chatId": %q}`, chatID))

	s.closeChatUseCase.EXPECT().Handle(eCtx.Request().Context(), closechat.Request{
		ID:        reqID,
		ManagerID: s.managerID,
		ChatID:    chatID,
	}).Return(closechat.ErrInvalidRequest)

	// Action.
	err := s.handlers.PostCloseChat(eCtx, managerv1.PostCloseChatParams{XRequestID: reqID})

	// Assert.
	s.Require().Error(err)
	s.Equal(http.StatusBadRequest, internalerrors.GetServerErrorCode(err))
	s.Empty(resp.Body)
}

func (s *HandlersSuite) TestCloseChat_UseCase_NoProblemError() {
	// Arrange.
	reqID := types.NewRequestID()
	chatID := types.NewChatID()

	resp, eCtx := s.newEchoCtx(reqID, "/v1/closeChat", fmt.Sprintf(`{"chatId": %q}`, chatID))

	s.closeChatUseCase.EXPECT().Handle(eCtx.Request().Context(), closechat.Request{
		ID:        reqID,
		ManagerID: s.managerID,
		ChatID:    chatID,
	}).Return(closechat.ErrProblemNotFound)

	// Action.
	err := s.handlers.PostCloseChat(eCtx, managerv1.PostCloseChatParams{XRequestID: reqID})

	// Assert.
	s.Require().Error(err)
	s.EqualValues(managerv1.ErrorCodeNoFoundProblem, internalerrors.GetServerErrorCode(err))
	s.Empty(resp.Body)
}

func (s *HandlersSuite) TestCloseChat_UseCase_UnknownError() {
	// Arrange.
	reqID := types.NewRequestID()
	chatID := types.NewChatID()

	resp, eCtx := s.newEchoCtx(reqID, "/v1/closeChat", fmt.Sprintf(`{"chatId": %q}`, chatID))

	s.closeChatUseCase.EXPECT().Handle(eCtx.Request().Context(), closechat.Request{
		ID:        reqID,
		ManagerID: s.managerID,
		ChatID:    chatID,
	}).Return(errors.New("unexpected"))

	// Action.
	err := s.handlers.PostCloseChat(eCtx, managerv1.PostCloseChatParams{XRequestID: reqID})

	// Assert.
	s.Require().Error(err)
	s.Equal(http.StatusInternalServerError, internalerrors.GetServerErrorCode(err))
	s.Empty(resp.Body)
}

func (s *HandlersSuite) TestCloseChat_UseCase_Success() {
	// Arrange.
	reqID := types.NewRequestID()
	chatID := types.NewChatID()

	resp, eCtx := s.newEchoCtx(reqID, "/v1/closeChat", fmt.Sprintf(`{"chatId": %q}`, chatID))

	s.closeChatUseCase.EXPECT().Handle(eCtx.Request().Context(), closechat.Request{
		ID:        reqID,
		ManagerID: s.managerID,
		ChatID:    chatID,
	}).Return(nil)

	// Action.
	err := s.handlers.PostCloseChat(eCtx, managerv1.PostCloseChatParams{XRequestID: reqID})

	// Assert.
	s.Require().NoError(err)
	s.Equal(http.StatusOK, resp.Code)
	s.JSONEq(`
{
   "data": null
}`, resp.Body.String())
}
