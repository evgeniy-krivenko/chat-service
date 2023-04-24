package errhandler

import (
	"fmt"
	"github.com/evgeniy-krivenko/chat-service/internal/errors"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
)

var _ echo.HTTPErrorHandler = Handler{}.Handle

//go:generate options-gen -out-filename=errhandler_options.gen.go -from-struct=Options
type Options struct {
	logger          *zap.Logger                                    `option:"mandatory" validate:"required"`
	productionMode  bool                                           `option:"mandatory"`
	responseBuilder func(code int, msg string, details string) any `option:"mandatory" validate:"required"`
}

type Handler struct {
	lg              *zap.Logger
	productionMode  bool
	responseBuilder func(code int, msg string, details string) any
}

func New(opts Options) (Handler, error) {
	if err := opts.Validate(); err != nil {
		return Handler{}, fmt.Errorf("validate errhandler options: %v", err)
	}
	return Handler{
		lg:              opts.logger,
		productionMode:  opts.productionMode,
		responseBuilder: opts.responseBuilder,
	}, nil
}

func (h Handler) Handle(err error, eCtx echo.Context) {
	code, msg, details := errors.ProcessServerError(err)
	if h.productionMode {
		details = ""
	}

	resp := h.responseBuilder(code, msg, details)
	err = eCtx.JSON(http.StatusOK, &resp)
	if err != nil {
		h.lg.Error(err.Error())
	}
}
