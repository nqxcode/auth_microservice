package cache

import (
	"context"

	"github.com/nqxcode/auth_microservice/internal/model"
)

func (s *service) GetMany(ctx context.Context, userIDs []int64) ([]model.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	panic("implement me")
}
