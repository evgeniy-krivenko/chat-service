package inmemmanagerpool

import (
	"context"
	"sync"

	"go.uber.org/zap"

	managerpool "github.com/evgeniy-krivenko/chat-service/internal/services/manager-pool"
	"github.com/evgeniy-krivenko/chat-service/internal/types"
)

const (
	serviceName = "manager-pool"
	managersMax = 1000
)

type Service struct {
	mu sync.RWMutex
	q  []types.UserID
	lg *zap.Logger
}

func (s *Service) Get(_ context.Context) (types.UserID, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.q) == 0 {
		return types.UserIDNil, managerpool.ErrNoAvailableManagers
	}

	var firstID types.UserID
	firstID, s.q = s.q[0], s.q[1:]

	s.lg.Info("manager removed", zap.Stringer("manager_id", firstID))
	return firstID, nil
}

func (s *Service) Put(ctx context.Context, managerID types.UserID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.contains(ctx, managerID) {
		return nil
	}

	s.q = append(s.q, managerID)

	s.lg.Info("manager stored", zap.Stringer("manager_id", managerID))
	return nil
}

func (s *Service) Contains(_ context.Context, managerID types.UserID) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, el := range s.q {
		if el.Matches(managerID) {
			return true, nil
		}
	}
	return false, nil
}

func (s *Service) Size() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.q)
}

func New() *Service {
	return &Service{
		q:  make([]types.UserID, 0, managersMax),
		lg: zap.L().Named(serviceName),
		mu: sync.RWMutex{},
	}
}

func (s *Service) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.q = s.q[0:]
	return nil
}

func (s *Service) contains(_ context.Context, managerID types.UserID) bool {
	for _, mID := range s.q {
		if mID == managerID {
			return true
		}
	}
	return false
}
