package user

import "context"

func (s *service) Delete(ctx context.Context, userID int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	keys, err := s.redisClient.Scan(ctx, buildListCacheKey("*"))
	if err != nil {
		return err
	}

	for _, key := range keys {
		err = s.redisClient.Delete(ctx, key)
		if err != nil {
			return err
		}
	}

	err = s.userRepository.Delete(ctx, userID)
	if err != nil {
		return err
	}

	return nil
}
