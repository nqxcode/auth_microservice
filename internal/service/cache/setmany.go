package cache

import (
	"context"

	"github.com/nqxcode/auth_microservice/internal/model"
)

func (s *service) SetList(ctx context.Context, users []model.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range users {
		_, err := s.userRepository.Create(ctx, &users[i])
		if err != nil {
			return err
		}
	}

	return nil
}
