package repository

import (
	"context"

	"github.com/nqxcode/auth_microservice/internal/model"
	"github.com/nqxcode/platform_common/pagination"
)

// UserRepository user repository
type UserRepository interface {
	Create(ctx context.Context, model *model.User) (int64, error)
	Update(ctx context.Context, id int64, info *model.UpdateUserInfo) error
	Delete(ctx context.Context, id int64) error
	Find(ctx context.Context, id int64) (*model.User, error)
	GetList(ctx context.Context, limit *pagination.Limit) ([]model.User, error)
}

// LogRepository log repository
type LogRepository interface {
	Create(ctx context.Context, model *model.Log) error
}
