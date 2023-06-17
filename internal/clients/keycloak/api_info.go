package keycloakclient

import (
	"fmt"

	"github.com/golang-jwt/jwt"

	"github.com/evgeniy-krivenko/chat-service/internal/types"
	"github.com/evgeniy-krivenko/chat-service/internal/validator"
)

type UserGetter struct{}

//nolint:tagliatelle
type claims struct {
	UserID     types.UserID `json:"sub,omitempty" validate:"required"`
	GivenName  string       `json:"given_name" validate:"omitempty"`
	FamilyName string       `json:"family_name" validate:"omitempty"`
}

func (c *claims) Valid() error {
	return validator.Validator.Struct(c)
}

type User struct {
	ID        types.UserID
	FirstName string
	LastName  string
}

func (c *UserGetter) GetUserInfoFromToken(tokenStr string) (*User, error) {
	var cl claims

	t, _, err := new(jwt.Parser).ParseUnverified(tokenStr, &cl)
	if err != nil {
		return nil, fmt.Errorf("parse claims: %v", err)
	}
	if err := t.Claims.Valid(); err != nil {
		return nil, fmt.Errorf("token claims not valid: %v", err)
	}

	return &User{
		ID:        cl.UserID,
		FirstName: cl.GivenName,
		LastName:  cl.FamilyName,
	}, nil
}
