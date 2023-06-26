package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	oapimdlwr "github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	keycloakclient "github.com/evgeniy-krivenko/chat-service/internal/clients/keycloak"
	"github.com/evgeniy-krivenko/chat-service/internal/middlewares"
)

const (
	readHeaderTimeout = time.Second
	shutdownTimeout   = 3 * time.Second
	bodySize          = "12KB"
)

type KeycloakClient interface {
	IntrospectToken(ctx context.Context, token string) (*keycloakclient.IntrospectTokenResult, error)
}

type wsHTTPHandler interface {
	Serve(eCtx echo.Context) error
}

//go:generate options-gen -out-filename=server_options.gen.go -from-struct=Options
type Options struct {
	logger           *zap.Logger           `option:"mandatory" validate:"required"`
	addr             string                `option:"mandatory" validate:"required,hostname_port"`
	allowOrigins     []string              `option:"mandatory" validate:"min=1"`
	v1Swagger        *openapi3.T           `option:"mandatory" validate:"required"`
	keycloakClient   KeycloakClient        `option:"mandatory" validate:"required"`
	resource         string                `option:"mandatory" validate:"required"`
	role             string                `option:"mandatory" validate:"required"`
	httpErrorHandler echo.HTTPErrorHandler `option:"mandatory" validate:"required"`
	registerHandlers func(*echo.Group)     `option:"mandatory" validate:"required"`
	shutdown         func()                `option:"mandatory" validate:"required"`

	wsHandler wsHTTPHandler `option:"mandatory" validate:"required"`
}

type Server struct {
	lg  *zap.Logger
	srv *http.Server
}

func New(opts Options) (*Server, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("validate server options: %v", err)
	}

	e := echo.New()

	e.GET(
		"/ws",
		opts.wsHandler.Serve,
		middlewares.NewKeyCloakWSTokenAuth(opts.keycloakClient, opts.resource, opts.role),
	)

	e.Use(
		middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: opts.allowOrigins,
			AllowMethods: []string{http.MethodPost},
		}),
		middlewares.NewLogger(opts.logger),
		middlewares.NewRecovery(opts.logger),
		middleware.BodyLimit(bodySize),
	)

	e.HTTPErrorHandler = opts.httpErrorHandler

	v1 := e.Group("v1",
		middlewares.NewKeyCloakTokenAuth(opts.keycloakClient, opts.resource, opts.role),
		oapimdlwr.OapiRequestValidatorWithOptions(opts.v1Swagger, &oapimdlwr.Options{
			Options: openapi3filter.Options{
				ExcludeRequestBody:  false,
				ExcludeResponseBody: true,
				AuthenticationFunc:  openapi3filter.NoopAuthenticationFunc,
			},
		}),
	)

	srv := &http.Server{
		Addr:              opts.addr,
		Handler:           e,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	srv.RegisterOnShutdown(opts.shutdown)

	opts.registerHandlers(v1)

	return &Server{
		lg:  opts.logger,
		srv: srv,
	}, nil
}

func (s *Server) Run(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(ctx, shutdownTimeout)
		defer cancel()

		return s.srv.Shutdown(ctx)
	})

	eg.Go(func() error {
		s.lg.Info("listen and serve", zap.String("addr", s.srv.Addr))

		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return fmt.Errorf("listen and serve: %v", err)
		}
		return nil
	})

	return eg.Wait()
}
