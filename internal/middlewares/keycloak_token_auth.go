package middlewares

import (
	"context"
	"errors"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	keycloakclient "github.com/evgeniy-krivenko/chat-service/internal/clients/keycloak"
	"github.com/evgeniy-krivenko/chat-service/internal/types"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/introspector_mock.gen.go -package=middlewaresmocks Introspector

const tokenCtxKey = "user-token"

var ErrNoRequiredResourceRole = errors.New("no required resource role")

type Introspector interface {
	IntrospectToken(ctx context.Context, token string) (*keycloakclient.IntrospectTokenResult, error)
}

// NewKeyCloakTokenAuth returns a middleware that implements "active" authentication:
// each request is verified by the Keycloak server.
func NewKeyCloakTokenAuth(introspector Introspector, resource, role string) echo.MiddlewareFunc {
	return middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup:  "header:Authorization",
		AuthScheme: "Bearer",
		Validator: func(tokenStr string, eCtx echo.Context) (bool, error) {
			ctx := eCtx.Request().Context()

			introspectResult, err := introspector.IntrospectToken(ctx, tokenStr)
			if err != nil {
				return false, err
			}

			if !introspectResult.Active {
				return false, nil
			}

			var cl claims

			token, _ := jwt.ParseWithClaims(
				tokenStr,
				&cl,
				func(token *jwt.Token) (interface{}, error) {
					return nil, nil
				},
			)

			if err := cl.Valid(); err != nil {
				return false, err
			}

			resource, ok := cl.ResourceAccess[resource]
			if !ok {
				return false, ErrNoRequiredResourceRole
			}

			roles, ok := resource["roles"]
			if !ok {
				return false, ErrNoRequiredResourceRole
			}

			for _, r := range roles {
				if r == role {
					eCtx.Set(tokenCtxKey, token)
					return true, nil
				}
			}
			return false, ErrNoRequiredResourceRole
		},
	})
}

func MustUserID(eCtx echo.Context) types.UserID {
	uid, ok := userID(eCtx)
	if !ok {
		panic("no user token in request context")
	}
	return uid
}

func userID(eCtx echo.Context) (types.UserID, bool) {
	t := eCtx.Get(tokenCtxKey)
	if t == nil {
		return types.UserIDNil, false
	}

	tt, ok := t.(*jwt.Token)
	if !ok {
		return types.UserIDNil, false
	}

	userIDProvider, ok := tt.Claims.(interface{ UserID() types.UserID })
	if !ok {
		return types.UserIDNil, false
	}
	return userIDProvider.UserID(), true
}
