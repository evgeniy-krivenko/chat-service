package managerload

import (
	"context"
	"fmt"

	"github.com/evgeniy-krivenko/chat-service/internal/types"
)

func (s *Service) CanManagerTakeProblem(ctx context.Context, managerID types.UserID) (bool, error) {
	count, err := s.problemsRepo.GetManagerOpenProblemsCount(ctx, managerID)
	if err != nil {
		return false, fmt.Errorf("get manager open problem: %v", err)
	}

	if count < s.maxProblemsAtTime {
		return true, nil
	}

	return false, nil
}
