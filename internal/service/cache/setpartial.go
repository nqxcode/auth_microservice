package cache

import (
	"context"

	"github.com/nqxcode/auth_microservice/internal/model"
)

func (s *service) SetPartial(ctx context.Context, id int64, user *model.UpdateUserInfo) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	err := s.userRepository.Update(ctx, id, user)
	if err != nil {
		return err
	}

	return nil
}
