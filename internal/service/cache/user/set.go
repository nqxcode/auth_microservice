package user

import (
	"context"

	"github.com/nqxcode/auth_microservice/internal/model"
)

// Set user to cache
func (s *service) Set(ctx context.Context, user *model.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.userRepository.Create(ctx, user)
	if err != nil {
		return err
	}

	err = s.invalidateLists(ctx)
	if err != nil {
		return err
	}

	return nil
}
