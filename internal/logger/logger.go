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

var atom zap.AtomicLevel

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
		return fmt.Errorf("validate logger options: %v", err)
	}

	atom = zap.NewAtomicLevel()

	level, err := zapcore.ParseLevel(opts.level)
	if err != nil {
		return fmt.Errorf("parsing lever: %v", err)
	}

	atom.SetLevel(level)

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
			atom,
		),
	}
	l := zap.New(zapcore.NewTee(cores...))
	zap.ReplaceGlobals(l)

	return nil
}

func SetLevel(l zapcore.Level) {
	atom.SetLevel(l)
}

func Sync() {
	if err := zap.L().Sync(); err != nil && !errors.Is(err, syscall.ENOTTY) {
		stdlog.Printf("cannot sync logger: %v", err)
	}
}
