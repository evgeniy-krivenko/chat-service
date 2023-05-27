package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/evgeniy-krivenko/chat-service/internal/buildinfo"
	keycloakclient "github.com/evgeniy-krivenko/chat-service/internal/clients/keycloak"
	"github.com/evgeniy-krivenko/chat-service/internal/config"
	"github.com/evgeniy-krivenko/chat-service/internal/logger"
	chatsrepo "github.com/evgeniy-krivenko/chat-service/internal/repositories/chats"
	jobsrepo "github.com/evgeniy-krivenko/chat-service/internal/repositories/jobs"
	messagesrepo "github.com/evgeniy-krivenko/chat-service/internal/repositories/messages"
	problemsrepo "github.com/evgeniy-krivenko/chat-service/internal/repositories/problems"
	clientevents "github.com/evgeniy-krivenko/chat-service/internal/server-client/events"
	clientv1 "github.com/evgeniy-krivenko/chat-service/internal/server-client/v1"
	serverdebug "github.com/evgeniy-krivenko/chat-service/internal/server-debug"
	managerv1 "github.com/evgeniy-krivenko/chat-service/internal/server-manager/v1"
	afcverdictsprocessor "github.com/evgeniy-krivenko/chat-service/internal/services/afc-verdicts-processor"
	inmemeventstream "github.com/evgeniy-krivenko/chat-service/internal/services/event-stream/in-mem"
	managerload "github.com/evgeniy-krivenko/chat-service/internal/services/manager-load"
	msgproducer "github.com/evgeniy-krivenko/chat-service/internal/services/msg-producer"
	"github.com/evgeniy-krivenko/chat-service/internal/services/outbox"
	clientmessageblockedjob "github.com/evgeniy-krivenko/chat-service/internal/services/outbox/jobs/client-message-blocked"
	clientmessagesentjob "github.com/evgeniy-krivenko/chat-service/internal/services/outbox/jobs/client-message-sent"
	sendclientmessagejob "github.com/evgeniy-krivenko/chat-service/internal/services/outbox/jobs/send-client-message"
	"github.com/evgeniy-krivenko/chat-service/internal/store"
	"github.com/evgeniy-krivenko/chat-service/internal/store/migrate"
	websocketstream "github.com/evgeniy-krivenko/chat-service/internal/websocket-stream"
)

const shutdownTimeout = 3 * time.Second

var configPath = flag.String("config", "configs/config.toml", "Path to config file")

func main() {
	if err := run(); err != nil {
		log.Fatalf("run app: %v", err)
	}
}

//nolint:gocyclo
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
		logger.WithDsnSentry(cfg.Sentry.DSN),
		logger.WithEnv(cfg.Global.Env),
	)

	err = logger.Init(logOpts)
	if err != nil {
		return fmt.Errorf("init logger with opts %+v: %v", logOpts, err)
	}
	defer logger.Sync()

	// Swaggers.
	swaggerClientV1, err := clientv1.GetSwagger()
	if err != nil {
		return fmt.Errorf("get client swagger: %v", err)
	}

	swaggerManagerV1, err := managerv1.GetSwagger()
	if err != nil {
		return fmt.Errorf("get manager swagger: %v", err)
	}

	swaggerEvents, err := clientevents.GetSwagger()
	if err != nil {
		return fmt.Errorf("get events swagger: %v", err)
	}

	// Debug server.
	srvDebug, err := serverdebug.New(serverdebug.NewOptions(
		cfg.Servers.Debug.Addr,
		swaggerClientV1,
		swaggerManagerV1,
		swaggerEvents,
	))
	if err != nil {
		return fmt.Errorf("init debug server: %v", err)
	}

	clientUserAgent := fmt.Sprintf("chat-service/%s", buildinfo.BuildInfo.Main.Version)

	// Database and migrations.
	storeClient, err := store.NewPSQLClient(store.NewPSQLOptions(
		cfg.DB.Address,
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Database,
		store.WithDebug(cfg.DB.DebugMode),
	))
	if err != nil {
		return fmt.Errorf("init psql client: %v", err)
	}
	defer storeClient.Close()

	// Migrations.
	if err := storeClient.Schema.Create(
		ctx,
		migrate.WithDropIndex(true),
		migrate.WithDropColumn(true),
	); err != nil {
		return fmt.Errorf("create migration: %v", err)
	}

	database := store.NewDatabase(storeClient)

	// Repositories.
	msgRepo, err := messagesrepo.New(messagesrepo.NewOptions(database))
	if err != nil {
		return fmt.Errorf("init messages repo: %v", err)
	}

	chatsRepo, err := chatsrepo.New(chatsrepo.NewOptions(database))
	if err != nil {
		return fmt.Errorf("init chats repo: %v", err)
	}

	problemsRepo, err := problemsrepo.New(problemsrepo.NewOptions(database))
	if err != nil {
		return fmt.Errorf("init problems repo: %v", err)
	}

	jobsRepo, err := jobsrepo.New(jobsrepo.NewOptions(database))
	if err != nil {
		return fmt.Errorf("init jobs repo: %v", err)
	}

	// Services.
	msgProducer, err := msgproducer.New(msgproducer.NewOptions(
		msgproducer.NewKafkaWriter(
			cfg.Services.MsgProducer.Brokers,
			cfg.Services.MsgProducer.Topic,
			cfg.Services.MsgProducer.BatchSize,
		),
		msgproducer.WithEncryptKey(cfg.Services.MsgProducer.EncryptKey),
	))
	if err != nil {
		return fmt.Errorf("init msg producer: %v", err)
	}

	outboxService, err := outbox.New(outbox.NewOptions(
		cfg.Services.Outbox.Workers,
		cfg.Services.Outbox.IDLE,
		cfg.Services.Outbox.ReserveFor,
		jobsRepo,
		database,
	))
	if err != nil {
		return fmt.Errorf("init outbox service: %v", err)
	}

	managerLoadService, err := managerload.New(managerload.NewOptions(
		cfg.Services.ManagerLoad.MaxProblemsAtSameTime,
		problemsRepo,
	))
	if err != nil {
		return fmt.Errorf("init manager load service: %v", err)
	}

	eventStream := inmemeventstream.New()
	defer eventStream.Close()

	afcVerdictProcessor, err := afcverdictsprocessor.New(afcverdictsprocessor.NewOptions(
		cfg.Services.AFCVerdictsProcessor.Brokers,
		cfg.Services.AFCVerdictsProcessor.Consumers,
		cfg.Services.AFCVerdictsProcessor.ConsumerGroup,
		cfg.Services.AFCVerdictsProcessor.VerdictsTopic,
		afcverdictsprocessor.NewKafkaReader,
		afcverdictsprocessor.NewKafkaDLQWriter(
			cfg.Services.AFCVerdictsProcessor.Brokers,
			cfg.Services.AFCVerdictsProcessor.VerdictsDLQTopic,
		),
		database,
		msgRepo,
		outboxService,
		afcverdictsprocessor.WithVerdictsSignKey(cfg.Services.AFCVerdictsProcessor.VerdictsSigningPublicKey),
	))
	if err != nil {
		return fmt.Errorf("init afc verdict processor: %v", err)
	}

	// Outbox Jobs.
	sendClientMessageJob, err := sendclientmessagejob.New(sendclientmessagejob.NewOptions(
		msgProducer,
		msgRepo,
		eventStream,
	))
	if err != nil {
		return fmt.Errorf("create send client message job: %v", err)
	}

	clientMessageSentJob, err := clientmessagesentjob.New(clientmessagesentjob.NewOptions(msgRepo, eventStream))
	if err != nil {
		return fmt.Errorf("create client msg sent job: %v", err)
	}

	clientMessageBlockedJob, err := clientmessageblockedjob.New(clientmessageblockedjob.NewOptions(msgRepo, eventStream))
	if err != nil {
		return fmt.Errorf("create client msg block job: %v", err)
	}

	// Register jobs
	outboxService.MustRegisterJob(sendClientMessageJob)
	outboxService.MustRegisterJob(clientMessageSentJob)
	outboxService.MustRegisterJob(clientMessageBlockedJob)

	// Clients.
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

	shutdownClient := make(chan struct{})
	shutdownManager := make(chan struct{})

	// Websocket client stream.
	wsClient, err := websocketstream.NewHTTPHandler(websocketstream.NewOptions(
		zap.L().Named("websocket-client"),
		eventStream,
		clientevents.Adapter{},
		websocketstream.JSONEventWriter{},
		websocketstream.NewUpgrader(cfg.Servers.Client.AllowOrigins, cfg.Servers.Client.SecWSProtocol),
		shutdownClient,
	))
	if err != nil {
		return fmt.Errorf("init websocket client: %v", err)
	}

	// Websocket manager stream.
	wsManager, err := websocketstream.NewHTTPHandler(websocketstream.NewOptions(
		zap.L().Named("websocket-manager"),
		eventStream,
		clientevents.Adapter{},
		websocketstream.JSONEventWriter{},
		websocketstream.NewUpgrader(cfg.Servers.Manager.AllowOrigins, cfg.Servers.Manager.SecWSProtocol),
		shutdownManager,
	))
	if err != nil {
		return fmt.Errorf("init websocket manager: %v", err)
	}

	// Servers.
	srvClient, err := initServerClient(
		cfg.Servers.Client.Addr,
		cfg.Servers.Client.AllowOrigins,
		swaggerClientV1,
		keycloakClient,
		wsClient,

		cfg.Servers.Client.RequiredAccess.Resource,
		cfg.Servers.Client.RequiredAccess.Role,

		msgRepo,
		chatsRepo,
		problemsRepo,
		outboxService,
		database,

		cfg.Global.IsProduction(),
	)
	if err != nil {
		return fmt.Errorf("init server client: %v", err)
	}

	srvManager, err := initServerManager(
		cfg.Servers.Manager.Addr,
		cfg.Servers.Manager.AllowOrigins,

		swaggerManagerV1,
		keycloakClient,
		wsManager,

		cfg.Servers.Manager.RequiredAccess.Resource,
		cfg.Servers.Manager.RequiredAccess.Role,

		managerLoadService,

		cfg.Global.IsProduction(),
	)
	if err != nil {
		return fmt.Errorf("init server manager: %v", err)
	}

	eg, ctx := errgroup.WithContext(ctx)

	// Run servers.
	eg.Go(func() error { return srvDebug.Run(ctx) })
	eg.Go(func() error { return srvClient.Run(ctx) })
	eg.Go(func() error { return srvManager.Run(ctx) })

	// Run services.
	eg.Go(func() error { return outboxService.Run(ctx) })
	eg.Go(func() error { return afcVerdictProcessor.Run(ctx) })

	// Websockets shutdown.
	eg.Go(func() error {
		<-ctx.Done()

		select {
		case shutdownClient <- struct{}{}:
		case <-time.After(shutdownTimeout):
		}

		return nil
	})

	eg.Go(func() error {
		<-ctx.Done()

		select {
		case shutdownManager <- struct{}{}:
		case <-time.After(shutdownTimeout):
		}

		return nil
	})

	if err = eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return fmt.Errorf("wait app stop: %v", err)
	}

	return nil
}
