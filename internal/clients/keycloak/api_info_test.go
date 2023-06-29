package keycloakclient_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	keycloakclient "github.com/evgeniy-krivenko/chat-service/internal/clients/keycloak"
)

func TestGetUserInfoFromToken(t *testing.T) {
	var usrGetter keycloakclient.UserGetter

	t.Run("success", func(t *testing.T) {
		const token = `eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpYXQiOjE2ODgwNjMyNzIsInN1YiI6ImQwZmZiZDM2LWJjMzAtMTFlZC04Mjg2LTQ2MWU0NjRlYmVkOCIsImdpdmVuX25hbWUiOiJFcmljIiwiZmFtaWx5X25hbWUiOiJDYXJ0bWFuIiwicmVzb3VyY2VfYWNjZXNzIjp7ImNoYXQtdWktbWFuYWdlciI6eyJyb2xlcyI6WyJzdXBwb3J0LWNoYXQtbWFuYWdlciJdfX19.zTjjgh6iiZQ-rHnE0kjDVqmdImlfDoE5TbxloM_deT8` //nolint:lll,gosec

		user, err := usrGetter.GetUserInfoFromToken(token)
		require.NoError(t, err)
		assert.Equal(t, "Eric", user.FirstName)
		assert.Equal(t, "Cartman", user.LastName)
	})

	t.Run("error", func(t *testing.T) {
		const token = `eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpYXQiOjE2ODY0MjA1NDgsInN1YiI6ImJjN2UzMDBiLTI5ZTQtNDdkNS1iYzkwLSIsImdpdmVuX25hbWUiOiJFcmljIiwiZmFtaWx5X25hbWUiOiJDYXJ0bWFuIn0.fkdx0NGM4aTRyHYrUY-m2LBqlLCqidVdnpRe_jih9ns` //nolint:lll,gosec

		user, err := usrGetter.GetUserInfoFromToken(token)
		require.Error(t, err)
		assert.Empty(t, user)
	})
}
