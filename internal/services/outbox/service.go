package outbox

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"

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
	lg       *zap.Logger
	Options
}

func New(opts Options) (*Service, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("validate outbox service options: %v", err)
	}
	return &Service{
		registry: make(map[string]Job),
		lg:       zap.L().Named(serviceName),
		Options:  opts,
	}, nil
}

func (s *Service) RegisterJob(job Job) error {
	if _, ok := s.registry[job.Name()]; ok {
		return errors.New("job exists")
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
	var wg sync.WaitGroup
	wg.Add(s.workers)
	for i := 0; i < s.workers; i++ {
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				default:
					s.worker(ctx)
				}
			}
		}()
	}
	s.lg.Info("service run", zap.Int("workers", s.workers), zap.Duration("idle time", s.idleTime))
	wg.Wait()
	return nil
}

func (s *Service) job(name string) (Job, bool) {
	j, ok := s.registry[name]
	return j, ok
}

func (s *Service) worker(ctx context.Context) {
	if err := s.work(ctx); err != nil {
		s.lg.Error("handle job", zap.Error(err))
	}
}

func (s *Service) work(ctx context.Context) error {
	jobInfo, err := s.jobsRepo.FindAndReserveJob(ctx, time.Now().Add(s.reserveFor))
	if errors.Is(err, jobsrepo.ErrNoJobs) {
		time.Sleep(s.idleTime)
		return nil
	}
	if err != nil {
		return fmt.Errorf("find and reserve job: %v", err)
	}

	job, ok := s.job(jobInfo.Name)
	if !ok {
		return s.moveToFailedJobWithReason(ctx, jobInfo, "no find registered job")
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, job.ExecutionTimeout())
	defer cancel()

	if err := job.Handle(ctxWithTimeout, jobInfo.Payload); err != nil {
		if jobInfo.Attempts >= job.MaxAttempts() {
			return s.moveToFailedJobWithReason(ctx, jobInfo, "max attempts exceeded")
		}

		return fmt.Errorf("job name %s, job id %s: %v", job.Name(), jobInfo.ID, err)
	}

	if err := s.jobsRepo.DeleteJob(ctx, jobInfo.ID); err != nil {
		return fmt.Errorf("delete job with ID %s: %v", jobInfo.ID.String(), err)
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

		s.lg.Warn(
			"job will be move to failed queue",
			zap.String("name", job.Name),
			zap.String("reason", reason),
		)

		return nil
	})
}
