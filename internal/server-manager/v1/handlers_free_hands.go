package managerv1

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/evgeniy-krivenko/chat-service/internal/errors"
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
	if err != nil {
		return errors.NewServerError[int](http.StatusInternalServerError, "internal error", err)
	}

	response := GetFreeHandsBtnAvailability{
		Available: useCaseResponse.Result,
	}

	return eCtx.JSON(http.StatusOK, &GetFreeHandsBtnAvailabilityResponse{Data: &response})
}
