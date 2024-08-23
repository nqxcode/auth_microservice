package cache

import (
	"context"
	"sync"

	"github.com/nqxcode/auth_microservice/internal/model"
)

func (s *service) SetList(ctx context.Context, users []model.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	errChan := make(chan error, len(users))

	var wg sync.WaitGroup
	for i := range users {
		wg.Add(1)
		go func(user model.User) {
			defer wg.Done()
			_, err := s.userRepository.Create(ctx, &user)
			if err != nil {
				errChan <- err
				return
			}
		}(users[i])
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}
