package keycloakclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

// StringOrArray represents a value that can either be a string or an array of strings.
type StringOrArray []string

func (s *StringOrArray) UnmarshalJSON(data []byte) error {
	if len(data) > 1 && data[0] == '[' {
		var obj []string
		if err := json.Unmarshal(data, &obj); err != nil {
			return err
		}
		*s = obj
		return nil
	}

	var obj string
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}
	*s = []string{obj}
	return nil
}

type IntrospectTokenResult struct {
	Exp    int           `json:"exp"`
	Iat    int           `json:"iat"`
	Aud    StringOrArray `json:"aud"`
	Active bool          `json:"active"`
}

// IntrospectToken implements
// https://www.keycloak.org/docs/latest/authorization_services/index.html#obtaining-information-about-an-rpt
func (c *Client) IntrospectToken(ctx context.Context, token string) (*IntrospectTokenResult, error) {
	url := fmt.Sprintf("realms/%s/protocol/openid-connect/token/introspect", c.realm)

	var introspectResult IntrospectTokenResult

	resp, err := c.auth(ctx).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(map[string]string{
			"token_type_hint": "requesting_party_token",
			"token":           token,
		}).
		SetResult(&introspectResult).
		Post(url)
	if err != nil {
		return nil, fmt.Errorf("send request to keycloak: %v", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("errored introspect keycloak response: %v", err)
	}

	return &introspectResult, nil
}

func (c *Client) auth(ctx context.Context) *resty.Request {
	return c.cli.R().
		SetBasicAuth(c.clientID, c.clientSecret).
		SetContext(ctx)
}
