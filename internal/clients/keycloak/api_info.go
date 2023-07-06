package keycloakclient

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt"

	"github.com/evgeniy-krivenko/chat-service/internal/types"
	"github.com/evgeniy-krivenko/chat-service/internal/validator"
)

var ErrNoAllowedResources = errors.New("no allowed resources")

type UserGetter struct{}

//nolint:tagliatelle
type claims struct {
	jwt.StandardClaims
	UserID          types.UserID   `json:"sub,omitempty" validate:"required"`
	Aud             StringOrArray  `json:"aud"`
	GivenName       string         `json:"given_name" validate:"omitempty"`
	FamilyName      string         `json:"family_name" validate:"omitempty"`
	ResourcesAccess resourceAccess `json:"resource_access"`
}

func (c *claims) Valid() error {
	if err := c.StandardClaims.Valid(); err != nil {
		return err
	}

	if len(c.ResourcesAccess) == 0 {
		return ErrNoAllowedResources
	}

	return validator.Validator.Struct(c)
}

type User struct {
	ID              types.UserID
	FirstName       string
	LastName        string
	ResourcesAccess resourceAccess
}

func (c *UserGetter) GetUserInfoFromToken(tokenStr string) (*User, error) {
	var cl claims

	t, _, err := new(jwt.Parser).ParseUnverified(tokenStr, &cl)
	if err != nil {
		return nil, fmt.Errorf("parse claims: %v", err)
	}
	if err := t.Claims.Valid(); err != nil {
		return nil, fmt.Errorf("token claims not valid: %w", err)
	}

	return &User{
		ID:              cl.UserID,
		FirstName:       cl.GivenName,
		LastName:        cl.FamilyName,
		ResourcesAccess: cl.ResourcesAccess,
	}, nil
}

type resourceAccess map[string]struct {
	Roles []string `json:"roles"`
}

func (ra resourceAccess) HasResourceRole(resource, role string) bool {
	access, ok := ra[resource]
	if !ok {
		return false
	}

	for _, r := range access.Roles {
		if r == role {
			return true
		}
	}
	return false
}
