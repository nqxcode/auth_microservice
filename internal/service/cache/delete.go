package cache

import "context"

func (s *service) Delete(ctx context.Context, userID int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	panic("implement me")
	// TODO implement
}
