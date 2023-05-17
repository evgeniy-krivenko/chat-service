package outbox

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	jobsrepo "github.com/evgeniy-krivenko/chat-service/internal/repositories/jobs"
	"github.com/evgeniy-krivenko/chat-service/internal/types"
)

const serviceName = "outbox"

type jobsRepository interface {
	CreateJob(ctx context.Context, name, payload string, availableAt time.Time) (types.JobID, error)
	FindAndReserveJob(ctx context.Context, until time.Time) (jobsrepo.Job, error)
	CreateFailedJob(ctx context.Context, name, payload, reason string) error
	DeleteJob(ctx context.Context, jobID types.JobID) error
}

type transactor interface {
	RunInTx(ctx context.Context, f func(context.Context) error) error
}

//go:generate options-gen -out-filename=service_options.gen.go -from-struct=Options
type Options struct {
	workers    int            `option:"mandatory" validate:"min=1,max=32"`
	idleTime   time.Duration  `option:"mandatory" validate:"min=100ms,max=10s"`
	reserveFor time.Duration  `option:"mandatory" validate:"min=1s,max=10m"`
	jobsRepo   jobsRepository `option:"mandatory" validate:"required"`
	trxtor     transactor     `option:"mandatory" validate:"required"`
}

type Service struct {
	registry map[string]Job
	Options
}

func New(opts Options) (*Service, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("validate outbox service options: %v", err)
	}
	return &Service{
		registry: make(map[string]Job),
		Options:  opts,
	}, nil
}

func (s *Service) RegisterJob(job Job) error {
	if _, ok := s.registry[job.Name()]; ok {
		return fmt.Errorf("job %v exists", job.Name())
	}

	s.registry[job.Name()] = job
	return nil
}

func (s *Service) MustRegisterJob(job Job) {
	if err := s.RegisterJob(job); err != nil {
		panic(err)
	}
}

func (s *Service) Run(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)

	for i := 0; i < s.workers; i++ {
		logger := zap.L().Named(serviceName).With(zap.Int("worker", i+1))
		eg.Go(func() error {
			for {
				if err := s.processAvailableJob(ctx, logger); err != nil {
					if ctx.Err() != nil {
						return nil
					}
					logger.Warn("job error", zap.Error(err))
					return err
				}

				select {
				case <-ctx.Done():
					return nil
				case <-time.After(s.idleTime):
				}
			}
		})
	}

	return eg.Wait()
}

func (s *Service) processAvailableJob(ctx context.Context, lg *zap.Logger) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		if err := s.work(ctx, lg); err != nil {
			if errors.Is(err, jobsrepo.ErrNoJobs) {
				lg.Debug("no jobs to process")
				return nil
			}

			return fmt.Errorf("process job: %v", err)
		}
	}
}

func (s *Service) job(name string) (Job, bool) {
	j, ok := s.registry[name]
	return j, ok
}

func (s *Service) work(ctx context.Context, lg *zap.Logger) error {
	jobInfo, err := s.jobsRepo.FindAndReserveJob(ctx, time.Now().Add(s.reserveFor))
	if err != nil {
		return fmt.Errorf("find and reserve job: %w", err)
	}

	job, ok := s.job(jobInfo.Name)
	if !ok {
		return s.moveToFailedJobWithReason(ctx, jobInfo, "no find registered job")
	}

	func() {
		ctx, cancel := context.WithTimeout(ctx, job.ExecutionTimeout())
		defer cancel()

		err = job.Handle(ctx, jobInfo.Payload)
	}()

	if err != nil {
		lg.Warn("hande job error", zap.Error(err))

		if jobInfo.Attempts >= job.MaxAttempts() {
			return s.moveToFailedJobWithReason(ctx, jobInfo, "max attempts exceeded")
		}

		return nil
	}

	if err := s.jobsRepo.DeleteJob(context.Background(), jobInfo.ID); err != nil {
		lg.Warn("delete job error", zap.Error(err))
	}

	return nil
}

func (s *Service) moveToFailedJobWithReason(ctx context.Context, job jobsrepo.Job, reason string) error {
	return s.trxtor.RunInTx(ctx, func(ctx context.Context) error {
		if err := s.jobsRepo.CreateFailedJob(ctx, job.Name, job.Payload, reason); err != nil {
			return fmt.Errorf("create failed job: %v", err)
		}

		if err := s.jobsRepo.DeleteJob(ctx, job.ID); err != nil {
			return fmt.Errorf("delete job while move to failed: %v", err)
		}

		return nil
	})
}
