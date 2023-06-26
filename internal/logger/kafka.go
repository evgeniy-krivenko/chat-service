package logger

import (
	"fmt"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type logType string

const (
	infoType  logType = "info"
	errorType logType = "error"
)

var _ kafka.Logger = (*KafkaAdapted)(nil)

type KafkaAdapted struct {
	lg      *zap.Logger
	logType logType
}

func (k KafkaAdapted) Printf(format string, v ...interface{}) {
	result := fmt.Sprintf(format, v...)
	switch k.logType {
	case infoType:
		k.lg.Debug(result) // for prevent more logs
	case errorType:
		k.lg.Error(result)
	}
}

func NewKafkaAdapted() *KafkaAdapted {
	return &KafkaAdapted{
		lg:      zap.L(),
		logType: infoType,
	}
}

func (k *KafkaAdapted) WithServiceName(name string) *KafkaAdapted {
	k.lg = k.lg.Named(name)
	return k
}

func (k *KafkaAdapted) ForErrors() *KafkaAdapted {
	k.logType = errorType
	return k
}
