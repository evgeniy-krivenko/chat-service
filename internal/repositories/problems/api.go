package problemsrepo

import (
	"context"
	"fmt"

	"github.com/evgeniy-krivenko/chat-service/internal/store"
	storeproblem "github.com/evgeniy-krivenko/chat-service/internal/store/problem"
	"github.com/evgeniy-krivenko/chat-service/internal/types"
)

func (r *Repo) CreateIfNotExists(ctx context.Context, chatID types.ChatID) (types.ProblemID, error) {
	existedProblemID, err := r.db.Problem(ctx).
		Query().
		Where(storeproblem.ChatIDEQ(chatID), storeproblem.ResolvedAtIsNil()).
		FirstID(ctx)
	if nil == err {
		return existedProblemID, nil
	}

	if !store.IsNotFound(err) {
		return types.ProblemIDNil, fmt.Errorf("select existent problem: %v", err)
	}

	newProblem, err := r.db.Problem(ctx).
		Create().
		SetChatID(chatID).
		Save(ctx)
	if err != nil {
		return types.ProblemIDNil, fmt.Errorf("create new problem: %v", err)
	}

	return newProblem.ID, nil
}

func (r *Repo) GetManagerOpenProblemsCount(ctx context.Context, managerID types.UserID) (int, error) {
	return r.db.Problem(ctx).Query().
		Where(
			storeproblem.ManagerID(managerID),
			storeproblem.ResolvedAtIsNil(),
		).
		Count(ctx)
}