package managerv1

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	internalerrors "github.com/evgeniy-krivenko/chat-service/internal/errors"
	"github.com/evgeniy-krivenko/chat-service/internal/middlewares"
	canreceiveproblems "github.com/evgeniy-krivenko/chat-service/internal/usecases/manager/can-receive-problems"
)

func (h Handlers) PostGetFreeHandsBtnAvailability(eCtx echo.Context, params PostGetFreeHandsBtnAvailabilityParams) error {
	ctx := eCtx.Request().Context()
	managerID := middlewares.MustUserID(eCtx)

	useCaseResponse, err := h.canReceiveProblemUseCase.Handle(ctx, canreceiveproblems.Request{
		ID:        params.XRequestID,
		ManagerID: managerID,
	})
	if errors.Is(err, canreceiveproblems.ErrInvalidRequest) {
		return internalerrors.NewServerError(http.StatusBadRequest, "invalid request", err)
	}
	if err != nil {
		return internalerrors.NewServerError(http.StatusInternalServerError, "internal error", err)
	}

	response := GetFreeHandsBtnAvailability{
		Available: useCaseResponse.Available,
		InPool:    useCaseResponse.InPool,
	}

	return eCtx.JSON(http.StatusOK, &GetFreeHandsBtnAvailabilityResponse{Data: &response})
}
