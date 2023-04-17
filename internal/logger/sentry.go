package logger

import (
	"crypto/x509"
	"fmt"
	"github.com/TheZeroSlave/zapsentry"
	"github.com/getsentry/sentry-go"
	"go.uber.org/zap/zapcore"
)

func NewSentryClient(dsn, env, version string) (*sentry.Client, error) {
	pool, err := x509.SystemCertPool()
	if err != nil {
		return nil, fmt.Errorf("load tsl certs: %v", err)
	}

	return sentry.NewClient(sentry.ClientOptions{
		Dsn:         dsn,
		Release:     version,
		Environment: env,
		CaCerts:     pool,
	})
}

func CoreSentry(client *sentry.Client) (zapcore.Core, error) {
	cfg := zapsentry.Configuration{
		Level: zapcore.WarnLevel,
	}
	return zapsentry.NewCore(cfg, zapsentry.NewSentryClientFromClient(client))
}
