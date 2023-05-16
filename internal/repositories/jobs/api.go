package jobsrepo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"

	"github.com/evgeniy-krivenko/chat-service/internal/store"
	storejob "github.com/evgeniy-krivenko/chat-service/internal/store/job"
	"github.com/evgeniy-krivenko/chat-service/internal/types"
)

var ErrNoJobs = errors.New("no jobs found")

type Job struct {
	ID       types.JobID
	Name     string
	Payload  string
	Attempts int
}

func (r *Repo) FindAndReserveJob(ctx context.Context, until time.Time) (Job, error) {
	var job *store.Job
	var err error

	err = r.db.RunInTx(ctx, func(ctx context.Context) error {
		job, err = r.db.Job(ctx).Query().
			ForUpdate(sql.WithLockAction(sql.SkipLocked)).
			Where(storejob.And(
				storejob.AvailableAtLTE(time.Now()),
				storejob.ReservedUntilLTE(time.Now()),
			)).
			Order(store.Asc(storejob.FieldCreatedAt)).
			First(ctx)

		if store.IsNotFound(err) {
			return ErrNoJobs
		}
		if err != nil {
			return err
		}

		job, err = job.Update().
			SetReservedUntil(until).
			AddAttempts(1).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("update job for reserve: %v", err)
		}

		return nil
	})
	if err != nil {
		return Job{}, fmt.Errorf("%w: find and reserve job", err)
	}

	return Job{
		ID:       job.ID,
		Name:     job.Name,
		Payload:  job.Payload,
		Attempts: job.Attempts,
	}, nil
}

func (r *Repo) CreateJob(ctx context.Context, name, payload string, availableAt time.Time) (types.JobID, error) {
	newJob, err := r.db.Job(ctx).Create().
		SetName(name).
		SetPayload(payload).
		SetReservedUntil(time.Now()).
		SetAvailableAt(availableAt).
		Save(ctx)
	if err != nil {
		return types.JobIDNil, err
	}
	return newJob.ID, nil
}

func (r *Repo) CreateFailedJob(ctx context.Context, name, payload, reason string) error {
	return r.db.FailedJob(ctx).Create().
		SetName(name).
		SetPayload(payload).
		SetReason(reason).
		Exec(ctx)
}

func (r *Repo) DeleteJob(ctx context.Context, jobID types.JobID) error {
	err := r.db.Job(ctx).DeleteOneID(jobID).Exec(ctx)
	if store.IsNotFound(err) {
		return ErrNoJobs
	}
	if err != nil {
		return err
	}

	return nil
}
