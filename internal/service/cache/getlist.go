package cache

import (
	"context"

	"github.com/nqxcode/auth_microservice/internal/model"
	"github.com/nqxcode/platform_common/pagination"
)

func (s *service) GetList(ctx context.Context, limit pagination.Limit) ([]model.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	users, err := s.userRepository.GetList(ctx, limit)
	if err != nil {
		return nil, err
	}

	return users, nil
}
