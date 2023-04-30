package middlewares

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/evgeniy-krivenko/chat-service/internal/errors"
)

func NewLogger(lg *zap.Logger) echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogLatency:   true,
		LogHost:      true,
		LogMethod:    true,
		LogRemoteIP:  true,
		LogRoutePath: true,
		LogRequestID: true,
		LogUserAgent: true,
		LogStatus:    true,
		LogError:     true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			status := v.Status
			if v.Error != nil {
				status = errors.GetServerErrorCode(v.Error)
			}

			fields := []zapcore.Field{
				zap.Duration("latency", v.Latency),
				zap.String("host", v.Host),
				zap.String("method", v.Method),
				zap.String("remote_id", v.RemoteIP),
				zap.String("path", v.RoutePath),
				zap.String("request_id", v.RequestID),
				zap.String("user_agent", v.UserAgent),
				zap.Int("status", status),
				zap.String("user_id", userIDForLog(c)),
				zap.Error(v.Error),
			}

			switch {
			case status >= 1000:
				lg.Error("business error", fields...)
			case status >= http.StatusInternalServerError:
				lg.Error("server error", fields...)
			case status >= http.StatusBadRequest:
				lg.Error("client err", fields...)
			default:
				lg.Info("success", fields...)
			}

			return nil
		},
	})
}

func userIDForLog(eCtx echo.Context) (id string) {
	uuid, ok := userID(eCtx)
	if ok && !uuid.IsZero() {
		id = uuid.String()
	}
	return
}
