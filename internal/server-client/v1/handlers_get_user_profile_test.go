package clientv1_test

import (
	"errors"
	"fmt"
	"net/http"

	internalerrors "github.com/evgeniy-krivenko/chat-service/internal/errors"
	clientv1 "github.com/evgeniy-krivenko/chat-service/internal/server-client/v1"
	"github.com/evgeniy-krivenko/chat-service/internal/types"
	getuserprofile "github.com/evgeniy-krivenko/chat-service/internal/usecases/client/get-user-profile"
)

func (s *HandlersSuite) TestGetUserProfile_UseCase_InvalidRequest() {
	// Arrange.
	reqID := types.NewRequestID()
	resp, eCtx := s.newEchoCtx(reqID, "/v1/getUserProfile", "")
	s.getUserProfileUseCase.EXPECT().Handle(eCtx.Request().Context(), getuserprofile.Request{
		ID:     reqID,
		UserID: s.clientID,
	}).Return(getuserprofile.Response{}, getuserprofile.ErrInvalidRequest)

	// Action.
	err := s.handlers.PostGetUserProfile(eCtx, clientv1.PostGetUserProfileParams{
		XRequestID: reqID,
	})

	// Assert.
	s.Require().Error(err)
	s.Equal(http.StatusBadRequest, internalerrors.GetServerErrorCode(err))
	s.Empty(resp.Body)
}

func (s *HandlersSuite) TestGetUserProfile_UseCase_ProfileNotFound() {
	// Arrange.
	reqID := types.NewRequestID()
	resp, eCtx := s.newEchoCtx(reqID, "/v1/getUserProfile", "")
	s.getUserProfileUseCase.EXPECT().Handle(eCtx.Request().Context(), getuserprofile.Request{
		ID:     reqID,
		UserID: s.clientID,
	}).Return(getuserprofile.Response{}, getuserprofile.ErrProfileNotFound)

	// Action.
	err := s.handlers.PostGetUserProfile(eCtx, clientv1.PostGetUserProfileParams{
		XRequestID: reqID,
	})

	// Assert.
	s.Require().Error(err)
	s.Equal(http.StatusBadRequest, internalerrors.GetServerErrorCode(err))
	s.Empty(resp.Body)
}

func (s *HandlersSuite) TestGetUserProfile_UseCase_UnknownError() {
	// Arrange.
	reqID := types.NewRequestID()
	resp, eCtx := s.newEchoCtx(reqID, "/v1/getUserProfile", "")
	s.getUserProfileUseCase.EXPECT().Handle(eCtx.Request().Context(), getuserprofile.Request{
		ID:     reqID,
		UserID: s.clientID,
	}).Return(getuserprofile.Response{}, errors.New("unexpected"))

	// Action.
	err := s.handlers.PostGetUserProfile(eCtx, clientv1.PostGetUserProfileParams{
		XRequestID: reqID,
	})

	// Assert.
	s.Require().Error(err)
	s.Equal(http.StatusInternalServerError, internalerrors.GetServerErrorCode(err))
	s.Empty(resp.Body)
}

func (s *HandlersSuite) TestGetUserProfile_UseCase_Success() {
	// Arrange.
	reqID := types.NewRequestID()
	firstName, lastName := "Eric", "Cartman"
	resp, eCtx := s.newEchoCtx(reqID, "/v1/getUserProfile", "")
	s.getUserProfileUseCase.EXPECT().Handle(eCtx.Request().Context(), getuserprofile.Request{
		ID:     reqID,
		UserID: s.clientID,
	}).Return(getuserprofile.Response{
		UserID:    s.clientID,
		FirstName: firstName,
		LastName:  lastName,
	}, nil)

	// Action.
	err := s.handlers.PostGetUserProfile(eCtx, clientv1.PostGetUserProfileParams{
		XRequestID: reqID,
	})

	// Assert.
	s.Require().NoError(err)
	s.JSONEq(fmt.Sprintf(`
	{
	  "data": {
		"id": %q,
		"firstName": %q,
		"lastName": %q
	  }
	}`, s.clientID, firstName, lastName), resp.Body.String())
}
