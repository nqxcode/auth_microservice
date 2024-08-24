package user

import (
	"context"

	"github.com/gomodule/redigo/redis"
	"github.com/nqxcode/auth_microservice/internal/model"
	"github.com/nqxcode/platform_common/helper/slice"
	"github.com/nqxcode/platform_common/pagination"
)

func (s *service) GetList(ctx context.Context, limit pagination.Limit) ([]model.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	ids, err := redis.Int64s(s.redisClient.LRange(ctx, buildListCacheKeyByLimit(limit), 0, -1))
	if err != nil {
		return nil, err
	}

	ids = slice.SliceByLimit(ids, limit)

	users, err := s.userRepository.GetByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	if len(users) != len(ids) {
		return nil, nil
	}

	return users, nil
}
