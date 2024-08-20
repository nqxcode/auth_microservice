package service

import (
	"context"

	"github.com/nqxcode/auth_microservice/internal/model"
)

// AuthService auth service
type AuthService interface {
	Create(ctx context.Context, user *model.User) (int64, error)
	Find(ctx context.Context, id int64) (*model.User, error)
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
