package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
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
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			lg.Info("success",
				zap.Duration("latency", v.Latency),
				zap.String("host", v.Host),
				zap.String("method", v.Method),
				zap.String("remote_id", v.RemoteIP),
				zap.String("path", v.RoutePath),
				zap.String("request_id", v.RequestID),
				zap.String("user_agent", v.UserAgent),
				zap.Int("status", v.Status),
				zap.String("user_id", userIDForLog(c)),
				zap.Error(v.Error),
			)

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
