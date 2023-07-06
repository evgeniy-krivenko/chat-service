package clientv1

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	internalerrors "github.com/evgeniy-krivenko/chat-service/internal/errors"
	"github.com/evgeniy-krivenko/chat-service/internal/middlewares"
	getuserprofile "github.com/evgeniy-krivenko/chat-service/internal/usecases/client/get-user-profile"
)

func (h Handlers) PostGetUserProfile(eCtx echo.Context, params PostGetUserProfileParams) error {
	ctx := eCtx.Request().Context()
	clientID := middlewares.MustUserID(eCtx)

	useCaseResponse, err := h.getUserProfile.Handle(ctx, getuserprofile.Request{
		ID:     params.XRequestID,
		UserID: clientID,
	})
	if err != nil {
		if errors.Is(err, getuserprofile.ErrInvalidRequest) {
			return internalerrors.NewServerError(
				http.StatusBadRequest,
				"invalid request for get client profile",
				err,
			)
		}
		if errors.Is(err, getuserprofile.ErrProfileNotFound) {
			return internalerrors.NewServerError(
				http.StatusBadRequest,
				"not found profile for client",
				err,
			)
		}

		return internalerrors.NewServerError(
			http.StatusInternalServerError,
			"unknown error while get client profile",
			err,
		)
	}

	return eCtx.JSON(http.StatusOK, &GetUserProfileResponse{Data: &UserProfile{
		Id:        useCaseResponse.UserID,
		FirstName: useCaseResponse.FirstName,
		LastName:  useCaseResponse.LastName,
	}})
}
