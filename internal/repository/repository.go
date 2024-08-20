package repository

import (
	"context"

	"github.com/nqxcode/auth_microservice/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, model *model.User) (int64, error)
	Update(ctx context.Context, id int64, info *model.UpdateUserInfo) error
	Delete(ctx context.Context, id int64) error
	//Find(ctx context.Context, id int64) (*model.User, error)
}

type LogRepository interface {
	Create(ctx context.Context, model *model.Log) error
}
