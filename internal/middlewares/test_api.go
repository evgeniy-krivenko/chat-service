package middlewares

import (
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"

	"github.com/evgeniy-krivenko/chat-service/internal/types"
)

func SetToken(c echo.Context, uid types.UserID) {
	t := jwt.NewWithClaims(jwt.SigningMethodRS256, &claims{
		StandardClaims: jwt.StandardClaims{},
		ResourceAccess: nil,
		userID:         uid,
	})
	c.Set(tokenCtxKey, t)
}
