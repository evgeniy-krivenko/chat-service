package problemsrepo

import (
	"context"
	"errors"
	"fmt"

	"github.com/evgeniy-krivenko/chat-service/internal/store"
	storemessage "github.com/evgeniy-krivenko/chat-service/internal/store/message"
	storeproblem "github.com/evgeniy-krivenko/chat-service/internal/store/problem"
	"github.com/evgeniy-krivenko/chat-service/internal/types"
	"github.com/evgeniy-krivenko/chat-service/pkg/utils"
)

var (
	ErrReqIDNotFount   = errors.New("request id not found")
	ErrProblemNotFound = errors.New("problem not found")
)

func (r *Repo) GetAvailableProblems(ctx context.Context) ([]Problem, error) {
	problems, err := r.db.Problem(ctx).Query().
		Where(storeproblem.ManagerIDIsNil()).
		Where(storeproblem.HasMessagesWith(
			storemessage.IsVisibleForManager(true),
		)).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("query problems: %v", err)
	}

	return utils.Apply(problems, adaptProblem), nil
}

func (r *Repo) SetManagerForProblem(ctx context.Context, problemID types.ProblemID, managerID types.UserID) error {
	_, err := r.db.Problem(ctx).
		UpdateOneID(problemID).
		Where(storeproblem.ManagerIDIsNil()).
		SetManagerID(managerID).
		Save(ctx)
	if err != nil {
		if store.IsNotFound(err) {
			return fmt.Errorf("set manager id %v for problem: %w", managerID, ErrProblemNotFound)
		}
		return fmt.Errorf("set manager id %v for problem: %v", managerID, err)
	}

	return nil
}

func (r *Repo) GetProblemRequestID(ctx context.Context, problemID types.ProblemID) (types.RequestID, error) {
	problem, err := r.db.Problem(ctx).Query().
		WithMessages(func(query *store.MessageQuery) {
			query.Where(storemessage.IsVisibleForManager(true)).
				Order(store.Asc(storemessage.FieldCreatedAt)).
				Limit(1)
		}).
		Where(storeproblem.ID(problemID)).
		Where(storeproblem.ManagerIDIsNil()).
		First(ctx)
	if err != nil {
		if store.IsNotFound(err) {
			return types.RequestIDNil, fmt.Errorf("problem id %v: %w", problemID, ErrReqIDNotFount)
		}

		return types.RequestIDNil, fmt.Errorf("problem id %v: %v", problemID, err)
	}

	if len(problem.Edges.Messages) > 0 {
		m := problem.Edges.Messages[0]
		return m.InitialRequestID, nil
	}

	return types.RequestIDNil, ErrReqIDNotFount
}
