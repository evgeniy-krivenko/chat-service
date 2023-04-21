package middlewares

import (
	"errors"

	"github.com/golang-jwt/jwt"

	"github.com/evgeniy-krivenko/chat-service/internal/types"
)

var (
	ErrNoAllowedResources = errors.New("no allowed resources")
	ErrSubjectNotDefined  = errors.New(`"sub" is not defined`)
)

type claims struct {
	jwt.StandardClaims
	ResourceAccess map[string]map[string][]string `json:"resource_access"` //nolint:tagliatelle
	userID         types.UserID
}

// Valid returns errors:
// - from StandardClaims validation;
// - ErrNoAllowedResources, if claims doesn't contain `resource_access` map or it's empty;
// - ErrSubjectNotDefined, if claims doesn't contain `sub` field or subject is zero UUID.
func (c *claims) Valid() error {
	if err := c.StandardClaims.Valid(); err != nil {
		return err
	}

	if c.ResourceAccess == nil || len(c.ResourceAccess) == 0 {
		return ErrNoAllowedResources
	}

	ui, err := types.Parse[types.UserID](c.Subject)
	if err != nil || ui.IsZero() {
		return ErrSubjectNotDefined
	}

	c.userID = ui

	return nil
}

func (c claims) UserID() types.UserID {
	return c.userID
}
