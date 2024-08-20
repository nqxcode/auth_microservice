package service

import (
	"context"
	"github.com/nqxcode/auth_microservice/internal/model"
)

type AuthService interface {
	Create(ctx context.Context, user *model.User) (int64, error)
	//Get(ctx context.Context, chat *model.User) (int64, error)
	Update(ctx context.Context, id int64, info *model.UpdateUserInfo) error
	Delete(ctx context.Context, id int64) error
}

type LogService interface {
	Create(ctx context.Context, message *model.Log) error
}

type HashService interface {
	Hash(ctx context.Context, password string) (string, error)
	GenerateSalt(ctx context.Context) (string, error)
}
