package problemsrepo

import (
	"context"
	"fmt"
	"time"

	storeproblem "github.com/evgeniy-krivenko/chat-service/internal/store/problem"
	"github.com/evgeniy-krivenko/chat-service/internal/types"
)

func (r *Repo) Resolve(
	ctx context.Context,
	reqID types.RequestID,
	managerID types.UserID,
	chatID types.ChatID,
) error {
	c, err := r.db.Problem(ctx).Update().
		Where(
			storeproblem.ManagerID(managerID),
			storeproblem.ChatID(chatID),
			storeproblem.ResolvedAtIsNil(),
		).
		SetResolvedAt(time.Now()).
		SetResolveRequestID(reqID).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("resolve problem: %v", err)
	}

	if c == 0 {
		return ErrProblemNotFound
	}

	return nil
}
