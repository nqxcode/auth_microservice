package user

import (
	"context"

	"github.com/nqxcode/auth_microservice/internal/model"
	"github.com/nqxcode/platform_common/helper/slice"
	modelCommon "github.com/nqxcode/platform_common/model"
	"github.com/nqxcode/platform_common/pagination"
)

func (s *service) SetList(ctx context.Context, users []model.User, limit pagination.Limit) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	err := s.redisClient.RPush(ctx, buildListCacheKeyByLimit(limit), slice.ToAnySlice(modelCommon.ExtractIDs(users)))
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
