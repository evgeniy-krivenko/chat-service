package clientv1

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	internalerrors "github.com/evgeniy-krivenko/chat-service/internal/errors"
	"github.com/evgeniy-krivenko/chat-service/internal/usecases/client/login"
)

func (h Handlers) PostLogin(eCtx echo.Context) error {
	ctx := eCtx.Request().Context()

	var req LoginRequest

	if err := eCtx.Bind(&req); err != nil {
		return internalerrors.NewServerError(http.StatusBadRequest, "bind login request", err)
	}

	resp, err := h.loginUseCase.Handle(ctx, login.Request{
		Login:    req.Login,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, login.ErrInvalidRequest) {
			return internalerrors.NewServerError(http.StatusBadRequest, "invalid login request", err)
		}
		if errors.Is(err, login.ErrAuthClient) {
			return internalerrors.NewServerError(http.StatusUnauthorized, "wrong login or password", err)
		}
		if errors.Is(err, login.ErrNoResourceAccess) {
			return internalerrors.NewServerError(http.StatusUnauthorized, "you have no access to this resource", err)
		}
		return internalerrors.NewServerError(http.StatusInternalServerError, "internal error", err)
	}

	return eCtx.JSON(http.StatusOK, &LoginResponse{
		Data: &LoginInfo{
			Token: resp.Token,
			User: UserProfile{
				Id:        resp.ClientID,
				FirstName: resp.FirstName,
				LastName:  resp.LastName,
			},
		},
	})
}
