package cache

import (
	"context"
	"github.com/nqxcode/auth_microservice/internal/model"
)

func (s *service) Get(ctx context.Context, userID int64) (*model.User, error) {
	panic("implement me")
}
