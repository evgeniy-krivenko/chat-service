package managerv1_test

import (
	"errors"
	"fmt"
	"net/http"

	internalerrors "github.com/evgeniy-krivenko/chat-service/internal/errors"
	managerv1 "github.com/evgeniy-krivenko/chat-service/internal/server-manager/v1"
	"github.com/evgeniy-krivenko/chat-service/internal/types"
	getmanagerprofile "github.com/evgeniy-krivenko/chat-service/internal/usecases/manager/get-manager-profile"
)

func (s *HandlersSuite) TestGetUserProfile_UseCase_InvalidRequest() {
	// Arrange.
	reqID := types.NewRequestID()
	resp, eCtx := s.newEchoCtx(reqID, "/v1/getUserProfile", "")
	s.getManagerProfileUseCase.EXPECT().Handle(eCtx.Request().Context(), getmanagerprofile.Request{
		ID:        reqID,
		ManagerID: s.managerID,
	}).Return(getmanagerprofile.Response{}, getmanagerprofile.ErrInvalidRequest)

	// Action.
	err := s.handlers.PostGetManagerProfile(eCtx, managerv1.PostGetManagerProfileParams{
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
	s.getManagerProfileUseCase.EXPECT().Handle(eCtx.Request().Context(), getmanagerprofile.Request{
		ID:        reqID,
		ManagerID: s.managerID,
	}).Return(getmanagerprofile.Response{}, getmanagerprofile.ErrProfileNotFound)

	// Action.
	err := s.handlers.PostGetManagerProfile(eCtx, managerv1.PostGetManagerProfileParams{
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
	s.getManagerProfileUseCase.EXPECT().Handle(eCtx.Request().Context(), getmanagerprofile.Request{
		ID:        reqID,
		ManagerID: s.managerID,
	}).Return(getmanagerprofile.Response{}, errors.New("unexpected"))

	// Action.
	err := s.handlers.PostGetManagerProfile(eCtx, managerv1.PostGetManagerProfileParams{
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
	s.getManagerProfileUseCase.EXPECT().Handle(eCtx.Request().Context(), getmanagerprofile.Request{
		ID:        reqID,
		ManagerID: s.managerID,
	}).Return(getmanagerprofile.Response{
		ManagerID: s.managerID,
		FirstName: firstName,
		LastName:  lastName,
	}, nil)

	// Action.
	err := s.handlers.PostGetManagerProfile(eCtx, managerv1.PostGetManagerProfileParams{
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
	}`, s.managerID, firstName, lastName), resp.Body.String())
}
