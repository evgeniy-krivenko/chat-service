package logger

import (
	"errors"
	"fmt"
	stdlog "log"
	"os"
	"syscall"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//go:generate options-gen -out-filename=logger_options.gen.go -from-struct=Options
type Options struct {
	level          string `option:"mandatory" validate:"required,oneof=debug info warn error"`
	productionMode bool
}

func MustInit(opts Options) {
	if err := Init(opts); err != nil {
		panic(err)
	}
}

func Init(opts Options) error {
	if err := opts.Validate(); err != nil {
		return fmt.Errorf("validate options: %v", err)
	}

	level, err := zap.ParseAtomicLevel(opts.level)
	if err != nil {
		return fmt.Errorf("parsing lever: %v", err)
	}

	var encoder zapcore.Encoder
	encoderCfg := zap.NewDevelopmentEncoderConfig()
	encoderCfg.NameKey = "component"

	if opts.productionMode {
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	} else {
		encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	}

	cores := []zapcore.Core{
		zapcore.NewCore(
			encoder,
			zapcore.Lock(os.Stdout),
			level,
		),
	}
	l := zap.New(zapcore.NewTee(cores...))
	zap.ReplaceGlobals(l)

	return nil
}

func Sync() {
	if err := zap.L().Sync(); err != nil && !errors.Is(err, syscall.ENOTTY) {
		stdlog.Printf("cannot sync logger: %v", err)
	}
}
