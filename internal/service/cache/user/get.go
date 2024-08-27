package user

import (
	"context"

	"github.com/nqxcode/auth_microservice/internal/model"
)

// Get user by id
func (s *service) Get(ctx context.Context, userID int64) (*model.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, err := s.userRepository.Get(ctx, userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}
