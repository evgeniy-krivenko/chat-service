package serverdebug

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/sync/errgroup"

	"github.com/evgeniy-krivenko/chat-service/internal/buildinfo"
	"github.com/evgeniy-krivenko/chat-service/internal/logger"
	"github.com/evgeniy-krivenko/chat-service/internal/middlewares"
)

const (
	readHeaderTimeout = time.Second
	shutdownTimeout   = 3 * time.Second
)

//go:generate options-gen -out-filename=server_options.gen.go -from-struct=Options
type Options struct {
	addr      string      `option:"mandatory" validate:"required,hostname_port"`
	v1Client  *openapi3.T `option:"mandatory" validate:"required"`
	v1Manager *openapi3.T `option:"mandatory" validate:"required"`
	events    *openapi3.T `option:"mandatory" validate:"required"`
}

type Server struct {
	lg        *zap.Logger
	srv       *http.Server
	v1Client  *openapi3.T
	v1Manager *openapi3.T
	events    *openapi3.T
}

func New(opts Options) (*Server, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("validate build server options: %v", err)
	}

	lg := zap.L().Named("server-debug")

	e := echo.New()
	e.Use(middleware.Recover())

	s := &Server{
		lg: lg,
		srv: &http.Server{
			Addr:              opts.addr,
			Handler:           e,
			ReadHeaderTimeout: readHeaderTimeout,
		},
		v1Client:  opts.v1Client,
		v1Manager: opts.v1Manager,
		events:    opts.events,
	}
	index := newIndexPage()

	e.Use(
		middlewares.NewLogger(lg),
		middlewares.NewRecovery(lg),
	)

	e.GET("/version", s.Version)
	e.PUT("/log/level", s.SetLogLevel)
	e.GET("/debug/error", s.Error)
	e.GET("/schema/client", s.SchemaClient)
	e.GET("/schema/manager", s.SchemaManager)
	e.GET("/schema/events", s.SchemaEvents)

	index.addPage("/version", "Get build information")
	index.addPage("/debug/pprof/", "Go to std profiler")
	index.addPage("/debug/pprof/profile?seconds=30", "Take half-min profile")
	index.addPage("/debug/error", "Debug Sentry error event")
	index.addPage("/schema/client", "Get client OpenAPI specification")
	index.addPage("/schema/manager", "Get client OpenAPI specification")
	index.addPage("/schema/events", "Get events OpenAPI specification")

	e.GET("/", index.handler)

	r := http.NewServeMux()
	r.Handle("/", e)
	r.HandleFunc("/debug/pprof/", pprof.Index)
	r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/debug/pprof/profile", pprof.Profile)
	r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	r.HandleFunc("/debug/pprof/trace", pprof.Trace)

	s.srv.Handler = r
	return s, nil
}

func (s *Server) Run(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		return s.srv.Shutdown(ctx)
	})

	eg.Go(func() error {
		s.lg.Info("listen and serve", zap.String("addr", s.srv.Addr))

		if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("listen and serve: %v", err)
		}
		return nil
	})

	return eg.Wait()
}

func (s *Server) Version(eCtx echo.Context) error {
	return eCtx.JSON(http.StatusOK, buildinfo.BuildInfo)
}

func (s *Server) SetLogLevel(eCtx echo.Context) error {
	level := eCtx.FormValue("level")

	l, err := zapcore.ParseLevel(level)
	if err != nil {
		return fmt.Errorf("parse level: %v", err)
	}

	logger.SetLevel(l)
	s.lg.Info("setting log level", zap.String("level", level))
	return nil
}

func (s *Server) Error(eCtx echo.Context) error {
	zap.L().Error("look at my in Sentry")
	return eCtx.String(http.StatusOK, "event was sent")
}

func (s *Server) SchemaClient(eCtx echo.Context) error {
	return eCtx.JSON(http.StatusOK, &s.v1Client)
}

func (s *Server) SchemaManager(eCtx echo.Context) error {
	return eCtx.JSON(http.StatusOK, &s.v1Manager)
}

func (s *Server) SchemaEvents(eCtx echo.Context) error {
	return eCtx.JSON(http.StatusOK, &s.events)
}
