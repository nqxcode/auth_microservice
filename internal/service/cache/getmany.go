package cache

import (
	"context"

	"github.com/nqxcode/auth_microservice/internal/model"
)

func (s *service) GetMany(ctx context.Context, userIDs []int64) ([]model.User, error) {
	panic("implement me")
}
