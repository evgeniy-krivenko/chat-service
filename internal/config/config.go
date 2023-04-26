package config

type Config struct {
	Global  GlobalConfig  `toml:"global"`
	Log     LogConfig     `toml:"log"`
	Servers ServersConfig `toml:"servers"`
	Sentry  SentryConfig  `toml:"sentry"`
	Clients ClientsConfig `toml:"clients"`
}

type GlobalConfig struct {
	Env string `toml:"env" validate:"required,oneof=dev stage prod"`
}

func (c GlobalConfig) IsProduction() bool {
	return c.Env == "prod"
}

type LogConfig struct {
	Level string `toml:"level" validate:"required,oneof=debug info warn error"`
}

type ServersConfig struct {
	Debug  DebugServerConfig  `toml:"debug"`
	Client ClientServerConfig `toml:"client"`
}

type DebugServerConfig struct {
	Addr string `toml:"addr" validate:"required,hostname_port"`
}

type ClientServerConfig struct {
	Add            string         `toml:"addr" validate:"hostname_port"`
	AllowOrigins   []string       `toml:"allow_origins" validate:"dive,url,min=1"`
	RequiredAccess RequiredAccess `toml:"required_access" validate:"required,dive"`
}

type RequiredAccess struct {
	Resource string `toml:"resource" validate:"required"`
	Role     string `toml:"role" validate:"required"`
}

type SentryConfig struct {
	DNS string `toml:"DNS" validate:"url"`
}

type ClientsConfig struct {
	Keycloak KeycloakConfig `toml:"keycloak"`
}

type KeycloakConfig struct {
	BasePath     string `toml:"base_path" validate:"url"`
	Realm        string `toml:"realm" validate:"required"`
	ClientID     string `toml:"client_id" validate:"required"` //nolint:tagliatelle
	ClientSecret string `toml:"client_secret" validate:"required"`
	DebugMode    bool   `toml:"debug_mode" validate:"boolean"`
}
