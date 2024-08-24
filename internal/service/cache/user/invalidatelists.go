package user

import "context"

func (s *service) invalidateLists(ctx context.Context) error {
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

	return nil
}
