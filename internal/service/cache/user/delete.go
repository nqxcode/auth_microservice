package user

import "context"

func (s *service) Delete(ctx context.Context, userID int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	err := s.userRepository.Delete(ctx, userID)
	if err != nil {
		return err
	}

	err = s.invalidateLists(ctx)
	if err != nil {
		return err
	}

	return nil
}
