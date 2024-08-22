package service

import (
	"context"

	"github.com/nqxcode/auth_microservice/internal/model"
	"github.com/nqxcode/platform_common/pagination"
)

// AuthService auth service
type AuthService interface {
	Create(ctx context.Context, user *model.User) (int64, error)
	Find(ctx context.Context, id int64) (*model.User, error)
	GetList(ctx context.Context, limit *pagination.Limit) ([]model.User, error)
	Update(ctx context.Context, id int64, info *model.UpdateUserInfo) error
	Delete(ctx context.Context, id int64) error
}

// LogService log service
type LogService interface {
	Create(ctx context.Context, message *model.Log) error
}

// HashService hash service
type HashService interface {
	Hash(ctx context.Context, password string) (string, error)
	GenerateSalt(ctx context.Context) (string, error)
}

// CacheService cache service
type CacheService interface {
	Set(ctx context.Context, user *model.User) error
	SetMany(ctx context.Context, users []model.User) error
	Get(ctx context.Context, userID int64) (*model.User, error)
	GetMany(ctx context.Context, userIDs []int64) ([]model.User, error)
	Delete(ctx context.Context, userID int64) error
}
