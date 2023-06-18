package clientv1_test

import (
	"errors"
	"fmt"
	"net/http"

	internalerrors "github.com/evgeniy-krivenko/chat-service/internal/errors"
	"github.com/evgeniy-krivenko/chat-service/internal/types"
	"github.com/evgeniy-krivenko/chat-service/internal/usecases/client/login"
)

const (
	log       = "user"
	password  = "123456"
	firstName = "Eric"
	lastName  = "Cartman"
	token     = "token"
)

func (s *HandlersSuite) TestLoginUseCase_InvalidRequest() {
	// Arrange.
	reqID := types.NewRequestID()
	body := fmt.Sprintf(`{"login": %q, "password": %q }`, log, password)

	resp, eCtx := s.newEchoCtx(reqID, "/v1/login", body)
	s.loginUseCase.EXPECT().Handle(eCtx.Request().Context(), login.Request{
		Login:    log,
		Password: password,
	}).Return(login.Response{}, login.ErrInvalidRequest)

	// Action.
	err := s.handlers.PostLogin(eCtx)

	// Assert.
	s.Require().Error(err)
	s.Equal(http.StatusBadRequest, internalerrors.GetServerErrorCode(err))
	s.Empty(resp.Body)
}

func (s *HandlersSuite) TestLoginUseCase_AuthClientError() {
	// Arrange.
	reqID := types.NewRequestID()
	body := fmt.Sprintf(`{"login": %q, "password": %q }`, log, password)

	resp, eCtx := s.newEchoCtx(reqID, "/v1/login", body)
	s.loginUseCase.EXPECT().Handle(eCtx.Request().Context(), login.Request{
		Login:    log,
		Password: password,
	}).Return(login.Response{}, login.ErrAuthClient)

	// Action.
	err := s.handlers.PostLogin(eCtx)

	// Assert.
	s.Require().Error(err)
	s.Equal(http.StatusUnauthorized, internalerrors.GetServerErrorCode(err))
	s.Empty(resp.Body)
}

func (s *HandlersSuite) TestLoginUseCase_UnexpectedError() {
	// Arrange.
	reqID := types.NewRequestID()

	body := fmt.Sprintf(`{"login": %q, "password": %q }`, log, password)

	resp, eCtx := s.newEchoCtx(reqID, "/v1/login", body)
	s.loginUseCase.EXPECT().Handle(eCtx.Request().Context(), login.Request{
		Login:    log,
		Password: password,
	}).Return(login.Response{}, errors.New("unexpected"))

	// Action.
	err := s.handlers.PostLogin(eCtx)

	// Assert.
	s.Require().Error(err)
	s.Equal(http.StatusInternalServerError, internalerrors.GetServerErrorCode(err))
	s.Empty(resp.Body)
}

func (s *HandlersSuite) TestLoginUseCase_Success() {
	// Arrange.
	reqID := types.NewRequestID()
	body := fmt.Sprintf(`{"login": %q, "password": %q }`, log, password)

	clientID := types.NewUserID()

	resp, eCtx := s.newEchoCtx(reqID, "/v1/login", body)
	s.loginUseCase.EXPECT().Handle(eCtx.Request().Context(), login.Request{
		Login:    log,
		Password: password,
	}).Return(login.Response{
		Token:     token,
		ClientID:  clientID,
		FirstName: firstName,
		LastName:  lastName,
	}, nil)

	// Action.
	err := s.handlers.PostLogin(eCtx)

	// Assert.
	s.Require().NoError(err)
	s.Equal(http.StatusOK, resp.Code)
	s.JSONEq(fmt.Sprintf(`
	{
		"data": {
			"token": %q,
			"user": {
				"id": %q,
				"firstName": %q,
				"lastName": %q
			}
		}
	}
`, token, clientID.String(), firstName, lastName), resp.Body.String())
}
