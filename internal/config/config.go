package config

import "time"

type Config struct {
	Global   GlobalConfig   `toml:"global"`
	Log      LogConfig      `toml:"log"`
	Servers  ServersConfig  `toml:"servers"`
	Sentry   SentryConfig   `toml:"sentry"`
	Clients  ClientsConfig  `toml:"clients"`
	DB       DBConfig       `toml:"db"`
	Services ServicesConfig `toml:"services"`
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
	Debug   DebugServerConfig   `toml:"debug"`
	Client  ClientServerConfig  `toml:"client"`
	Manager ManagerServerConfig `toml:"manager"`
}

type DebugServerConfig struct {
	Addr string `toml:"addr" validate:"required,hostname_port"`
}

type ClientServerConfig struct {
	Addr           string         `toml:"addr" validate:"hostname_port"`
	AllowOrigins   []string       `toml:"allow_origins" validate:"dive,url,min=1"`
	SecWSProtocol  string         `toml:"sec_ws_protocol" validate:"required"`
	RequiredAccess RequiredAccess `toml:"required_access" validate:"required,dive"`
}

type ManagerServerConfig struct {
	Addr           string         `toml:"addr" validate:"hostname_port"`
	AllowOrigins   []string       `toml:"allow_origins" validate:"dive,url,min=1"`
	SecWSProtocol  string         `toml:"sec_ws_protocol" validate:"required"`
	RequiredAccess RequiredAccess `toml:"required_access" validate:"required,dive"`
}

type RequiredAccess struct {
	Resource string `toml:"resource" validate:"required"`
	Role     string `toml:"role" validate:"required"`
}

type SentryConfig struct {
	DSN string `toml:"dsn" validate:"url"`
}

type ClientsConfig struct {
	Keycloak KeycloakConfig `toml:"keycloak"`
}

type KeycloakConfig struct {
	BasePath     string `toml:"base_path" validate:"url"`
	Realm        string `toml:"realm" validate:"required"`
	ClientID     string `toml:"client_id" validate:"required"`
	ClientSecret string `toml:"client_secret" validate:"required"`
	DebugMode    bool   `toml:"debug_mode" validate:"boolean"`
}

type DBConfig struct {
	Address   string `toml:"addr" validate:"required,hostname_port"`
	Database  string `toml:"database" validate:"required"`
	User      string `toml:"user" validate:"required"`
	Password  string `toml:"password" validate:"required"`
	DebugMode bool   `toml:"debug_mode"`
}

type ServicesConfig struct {
	Outbox      OutboxConfig      `toml:"outbox"`
	MsgProducer MsgProducerConfig `toml:"msg_producer"`
	ManagerLoad ManagerLoadConfig `toml:"manager_load"`
}

type OutboxConfig struct {
	Workers    int           `toml:"workers" validate:"required,min=1"`
	IDLE       time.Duration `toml:"idle_time" validate:"required,min=500ms,max=10s"`
	ReserveFor time.Duration `toml:"reserve_for" validate:"required"`
}

type MsgProducerConfig struct {
	Brokers    []string `toml:"brokers" validate:"required,gt=0,dive,required,hostname_port"`
	Topic      string   `toml:"topic" validate:"required"`
	BatchSize  int      `toml:"batch_size" validate:"required,min=1"`
	EncryptKey string   `toml:"encrypt_key" validate:"omitempty,hexadecimal"`
}

type ManagerLoadConfig struct {
	MaxProblemsAtSameTime int `toml:"max_problems_at_same_time" validate:"required,gt=0"`
}
