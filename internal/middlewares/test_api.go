package middlewares

import (
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"

	"github.com/evgeniy-krivenko/chat-service/internal/types"
)

func AuthWith(uid types.UserID) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			SetToken(c, uid)
			return next(c)
		}
	}
}

func SetToken(c echo.Context, uid types.UserID) {
	t := jwt.NewWithClaims(jwt.SigningMethodRS256, &claims{
		StandardClaims: jwt.StandardClaims{},
		ResourceAccess: nil,
		userID:         uid,
	})
	c.Set(tokenCtxKey, t)
}
