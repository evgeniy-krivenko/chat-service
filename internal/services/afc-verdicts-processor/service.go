package afcverdictsprocessor

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/protocol"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	clientmessageblockedjob "github.com/evgeniy-krivenko/chat-service/internal/services/outbox/jobs/client-message-blocked"
	clientmessagesentjob "github.com/evgeniy-krivenko/chat-service/internal/services/outbox/jobs/client-message-sent"
	"github.com/evgeniy-krivenko/chat-service/internal/services/outbox/jobs/payload/simpleid"
	"github.com/evgeniy-krivenko/chat-service/internal/types"
)

const serviceName = "afc-verdicts-processor"

//go:generate mockgen -source=$GOFILE -destination=mocks/service_mock.gen.go -package=afcverdictsprocessormocks

type DLQFunc func(context.Context, ...kafka.Message) error

type TxFunc func(ctx context.Context, f func(context.Context) error) error

type messagesRepository interface {
	MarkAsVisibleForManager(ctx context.Context, msgID types.MessageID) error
	BlockMessage(ctx context.Context, msgID types.MessageID) error
}

type outboxService interface {
	Put(ctx context.Context, name, payload string, availableAt time.Time) (types.JobID, error)
}

type transactor interface {
	RunInTx(ctx context.Context, f func(context.Context) error) error
}

//go:generate options-gen -out-filename=service_options.gen.go -from-struct=Options
type Options struct {
	backoffInitialInterval time.Duration `default:"100ms" validate:"min=50ms,max=1s"`
	backoffMaxElapsedTime  time.Duration `default:"5s" validate:"min=500ms,max=1m"`
	expFactor              float64       `default:"2.71828" validate:"min=1.5,max=5.0"`
	expJitter              float64       `default:"0.1" validate:"min=0.1"`

	brokers          []string `option:"mandatory" validate:"min=1"`
	consumers        int      `option:"mandatory" validate:"min=1,max=16"`
	consumerGroup    string   `option:"mandatory" validate:"required"`
	verdictsTopic    string   `option:"mandatory" validate:"required"`
	verdictsSignKey  string
	processBatchSize int `default:"1"`

	readerFactory KafkaReaderFactory `option:"mandatory" validate:"required"`
	dlqWriter     KafkaDLQWriter     `option:"mandatory" validate:"required"`

	txtor   transactor         `option:"mandatory" validate:"required"`
	msgRepo messagesRepository `option:"mandatory" validate:"required"`
	outBox  outboxService      `option:"mandatory" validate:"required"`
}

type Service struct {
	Options
	lg        *zap.Logger
	publicKey *rsa.PublicKey
}

func New(opts Options) (*Service, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("validate afc verdicts processor options: %v", err)
	}

	s := &Service{
		Options: opts,
		lg:      zap.L().Named(serviceName),
	}

	var err error
	if k := opts.verdictsSignKey; k != "" {
		if s.publicKey, err = jwt.ParseRSAPublicKeyFromPEM([]byte(k)); err != nil {
			return nil, fmt.Errorf("parse public key afc verdicts: %v", err)
		}
	}

	return s, nil
}

func (s *Service) Run(ctx context.Context) error {
	defer s.dlqWriter.Close()

	s.lg.Info(
		"service run",
		zap.Int("consumers", s.consumers),
		zap.String("verdict topic", s.verdictsTopic),
	)

	eg, ctx := errgroup.WithContext(ctx)

	for i := 0; i < s.consumers; i++ {
		eg.Go(func() error { return s.consume(ctx) })
	}

	if err := eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		s.lg.Error("consume", zap.Error(err))

		return fmt.Errorf("wait for afc verdicts consumers: %v", err)
	}
	return nil
}

func (s *Service) consume(ctx context.Context) error {
	reader := s.readerFactory(s.brokers, s.consumerGroup, s.verdictsTopic)
	msgs := make([]kafka.Message, 0, s.processBatchSize)

	defer func() {
		// commit if msgs not empty when context cancel
		if len(msgs) > 0 {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
			if err := reader.CommitMessages(ctx, msgs...); err != nil {
				s.lg.Error("commit messages with ctx cancel")
			}
			cancel()
		}

		reader.Close()
	}()

	for {
		msg, err := reader.FetchMessage(ctx)
		if errors.Is(err, io.EOF) {
			return nil
		}
		if errors.Is(err, context.Canceled) {
			return nil
		}
		if err != nil {
			return fmt.Errorf("%w: fetch msg", err)
		}

		if err := s.handleMsg(ctx, msg); err != nil {
			s.lg.Error("handle msg", zap.Error(err))
			return fmt.Errorf("handle msg: %v", err)
		}

		msgs = append(msgs, msg)

		if len(msgs) == s.processBatchSize {
			if err := reader.CommitMessages(ctx, msgs...); err != nil {
				return fmt.Errorf("%w: commit msg", err)
			}

			msgs = msgs[:0]
		}
	}
}

func (s *Service) handleMsg(ctx context.Context, msg kafka.Message) error {
	var v Verdict
	verdictData := msg.Value

	parts := strings.Split(string(verdictData), ".")
	if len(parts) == 3 {
		if err := jwt.SigningMethodRS256.Verify(strings.Join(parts[0:2], "."), parts[2], s.publicKey); err != nil {
			s.writeToDLQ(ctx, msg, err.Error())
			return nil
		}

		data, err := jwt.DecodeSegment(parts[1])
		if err != nil {
			s.writeToDLQ(ctx, msg, err.Error())
			return nil
		}

		verdictData = data
	}

	if err := json.Unmarshal(verdictData, &v); err != nil {
		s.writeToDLQ(ctx, msg, err.Error())
		return nil
	}

	if err := v.Valid(); err != nil {
		s.writeToDLQ(ctx, msg, err.Error())
		return nil
	}

	if !v.IsSuccess() {
		if err := s.blockMessage(ctx, v); err != nil {
			s.writeToDLQ(ctx, msg, err.Error())
		}

		return nil
	}

	if err := s.markAsVisibleForManager(ctx, v); err != nil {
		s.writeToDLQ(ctx, msg, err.Error())
	}

	return nil
}

func (s *Service) blockMessage(ctx context.Context, v Verdict) error {
	txWithBackoff := s.backoffTx(s.txtor.RunInTx)

	return txWithBackoff(ctx, func(ctx context.Context) error {
		if err := s.msgRepo.BlockMessage(ctx, v.MessageID); err != nil {
			return fmt.Errorf("block msg: %v", err)
		}

		if _, err := s.outBox.Put(
			ctx,
			clientmessageblockedjob.Name,
			simpleid.MustMarshal(v.MessageID),
			time.Now(),
		); err != nil {
			return fmt.Errorf("put job %v to outbox: %v", clientmessageblockedjob.Name, err)
		}

		s.lg.Info("block suspicious message", zap.Stringer("message_id", v.MessageID))
		return nil
	})
}

func (s *Service) markAsVisibleForManager(ctx context.Context, v Verdict) error {
	txWithBackoff := s.backoffTx(s.txtor.RunInTx)

	return txWithBackoff(ctx, func(ctx context.Context) error {
		if err := s.msgRepo.MarkAsVisibleForManager(ctx, v.MessageID); err != nil {
			return fmt.Errorf("mark msg as visible for manager: %v", err)
		}

		if _, err := s.outBox.Put(
			ctx,
			clientmessagesentjob.Name,
			simpleid.MustMarshal(v.MessageID),
			time.Now(),
		); err != nil {
			return fmt.Errorf("put job %v to outbox: %v", clientmessagesentjob.Name, err)
		}

		s.lg.Info("mark as visible for manager", zap.Stringer("message_id", v.MessageID))
		return nil
	})
}

func (s *Service) writeToDLQ(ctx context.Context, msg kafka.Message, errMsg string) {
	lastErr := protocol.Header{
		Key:   "LAST_ERROR",
		Value: []byte(errMsg),
	}

	originPartition := protocol.Header{
		Key:   "ORIGINAL_PARTITION",
		Value: []byte(strconv.Itoa(msg.Partition)),
	}

	msg.Headers = append(msg.Headers, lastErr, originPartition)
	msg.Topic = ""

	go func() {
		if err := s.dlqWriter.WriteMessages(ctx, msg); err != nil {
			s.lg.Error("write msg to dlq after exponential backoff", zap.Error(err))
			return
		}
		s.lg.Warn("write msg to dlq")
	}()
}

func (s *Service) backoffTx(f TxFunc) TxFunc {
	return func(ctx context.Context, f2 func(context.Context) error) error {
		delay := s.backoffInitialInterval

		for {
			if err := f(ctx, f2); nil == err || delay >= s.backoffMaxElapsedTime {
				return err
			}

			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return nil
			}

			delay = time.Duration(float64(delay) * s.expFactor)
			if delay > s.backoffMaxElapsedTime {
				delay = s.backoffMaxElapsedTime
			}
			delay += time.Duration(rand.NormFloat64() * s.expJitter * float64(time.Second))
		}
	}
}

func (s *Service) backoffDLQ(f DLQFunc) DLQFunc {
	return func(ctx context.Context, msgs ...kafka.Message) error {
		delay := s.backoffInitialInterval

		for {
			if err := f(ctx, msgs...); nil == err || delay >= s.backoffMaxElapsedTime {
				return err
			}

			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return nil
			}

			delay = time.Duration(float64(delay) * s.expFactor)
			if delay > s.backoffMaxElapsedTime {
				delay = s.backoffMaxElapsedTime
			}
			delay += time.Duration(rand.NormFloat64() * s.expJitter * float64(time.Second))
		}
	}
}
