package user

import (
	"context"

	"github.com/nqxcode/auth_microservice/internal/model"
	modelCommon "github.com/nqxcode/platform_common/model"
	"github.com/nqxcode/platform_common/pagination"
)

func (s *service) SetList(ctx context.Context, users []model.User, limit pagination.Limit) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	err := s.redisClient.RPush(ctx, buildListCacheKeyByLimit(limit), func() []interface{} {
		userIDs := modelCommon.ExtractIDs(users)

		var results []interface{}
		for _, user := range userIDs {
			results = append(results, user)
		}
		return results
	}())
	if err != nil {
		return err
	}

	for _, user := range users {
		_, err = s.userRepository.Create(ctx, &user)
		if err != nil {
			return err
		}
	}

	return nil
}
