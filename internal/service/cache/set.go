package cache

import (
	"context"

	"github.com/nqxcode/auth_microservice/internal/model"
)

func (s *service) Set(ctx context.Context, user *model.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	panic("implement me")
}
