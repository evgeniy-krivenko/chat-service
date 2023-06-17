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
		const token = `eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpYXQiOjE2ODY0MTk5MTUsInN1YiI6ImJjN2UzMDBiLTI5ZTQtNDdkNS1iYzkwLTkwY2EwMDQ2ZjlmNyIsImdpdmVuX25hbWUiOiJFcmljIiwiZmFtaWx5X25hbWUiOiJDYXJ0bWFuIn0.Hjdr_26PAABM6Bnw_p8rMyCgbthdiGngu4OQVsCiwEk` //nolint:lll,gosec

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
