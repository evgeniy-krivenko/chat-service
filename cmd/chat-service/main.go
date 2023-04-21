package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"

	"github.com/evgeniy-krivenko/chat-service/internal/buildinfo"
	keycloakclient "github.com/evgeniy-krivenko/chat-service/internal/clients/keycloak"
	"github.com/evgeniy-krivenko/chat-service/internal/config"
	"github.com/evgeniy-krivenko/chat-service/internal/logger"
	clientv1 "github.com/evgeniy-krivenko/chat-service/internal/server-client/v1"
	serverdebug "github.com/evgeniy-krivenko/chat-service/internal/server-debug"
)

var configPath = flag.String("config", "configs/config.toml", "Path to config file")

func main() {
	if err := run(); err != nil {
		log.Fatalf("run app: %v", err)
	}
}

func run() (errReturned error) {
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cfg, err := config.ParseAndValidate(*configPath)
	if err != nil {
		return fmt.Errorf("parse and validate config %q: %v", *configPath, err)
	}

	logOpts := logger.NewOptions(
		cfg.Log.Level,
		logger.WithProductionMode(cfg.Global.IsProduction()),
		logger.WithDns(cfg.Sentry.DNS),
		logger.WithEnv(cfg.Global.Env),
	)

	err = logger.Init(logOpts)
	if err != nil {
		return fmt.Errorf("init logger with opts %+v: %v", logOpts, err)
	}
	defer logger.Sync()

	swaggerClientV1, err := clientv1.GetSwagger()
	if err != nil {
		return fmt.Errorf("get client swagger: %v", err)
	}

	srvDebug, err := serverdebug.New(serverdebug.NewOptions(
		cfg.Servers.Debug.Addr,
		swaggerClientV1,
	))
	if err != nil {
		return fmt.Errorf("init debug server: %v", err)
	}

	clientUserAgent := fmt.Sprintf("chat-service/%s", buildinfo.BuildInfo.Main.Version)

	keycloakClient, err := keycloakclient.New(keycloakclient.NewOptions(
		cfg.Clients.Keycloak.BasePath,
		cfg.Clients.Keycloak.Realm,
		cfg.Clients.Keycloak.ClientID,
		cfg.Clients.Keycloak.ClientSecret,
		keycloakclient.WithDebugMode(cfg.Clients.Keycloak.DebugMode),
		keycloakclient.WithUserAgent(clientUserAgent),
		keycloakclient.WithProductionMode(cfg.Global.IsProduction()),
	))
	if err != nil {
		return fmt.Errorf("init keycloak client: %v", err)
	}

	srvClient, err := initServerClient(
		cfg.Servers.Client.Add,
		cfg.Servers.Client.AllowOrigins,
		swaggerClientV1,
		keycloakClient,
		cfg.Servers.Client.RequiredAccess.Resource,
		cfg.Servers.Client.RequiredAccess.Role,
	)
	if err != nil {
		return fmt.Errorf("init server client: %v", err)
	}

	eg, ctx := errgroup.WithContext(ctx)

	// Run servers.
	eg.Go(func() error { return srvDebug.Run(ctx) })

	eg.Go(func() error { return srvClient.Run(ctx) })
	// Run services.
	// Ждут своего часа.
	// ...

	if err = eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return fmt.Errorf("wait app stop: %v", err)
	}

	return nil
}
