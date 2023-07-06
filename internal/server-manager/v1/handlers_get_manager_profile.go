package managerv1

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	internalerrors "github.com/evgeniy-krivenko/chat-service/internal/errors"
	"github.com/evgeniy-krivenko/chat-service/internal/middlewares"
	getmanagerprofile "github.com/evgeniy-krivenko/chat-service/internal/usecases/manager/get-manager-profile"
)

func (h Handlers) PostGetManagerProfile(eCtx echo.Context, params PostGetManagerProfileParams) error {
	ctx := eCtx.Request().Context()
	managerID := middlewares.MustUserID(eCtx)

	resp, err := h.getManagerProfile.Handle(ctx, getmanagerprofile.Request{
		ID:        params.XRequestID,
		ManagerID: managerID,
	})
	if err != nil {
		if errors.Is(err, getmanagerprofile.ErrInvalidRequest) {
			return internalerrors.NewServerError(
				http.StatusBadRequest,
				"invalid request for get manager profile",
				err,
			)
		}
		if errors.Is(err, getmanagerprofile.ErrProfileNotFound) {
			return internalerrors.NewServerError(
				http.StatusBadRequest,
				"not found profile for manager",
				err,
			)
		}

		return internalerrors.NewServerError(
			http.StatusInternalServerError,
			"unknown error while get manager profile",
			err,
		)
	}

	return eCtx.JSON(http.StatusOK, &GetManagerProfileResponse{Data: &ManagerProfile{
		Id:        resp.ManagerID,
		FirstName: resp.FirstName,
		LastName:  resp.LastName,
	}})
}
