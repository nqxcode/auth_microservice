package service

import (
	"context"
	"time"

	"github.com/nqxcode/auth_microservice/internal/model"
	"github.com/nqxcode/platform_common/pagination"
)

// AuthService auth service
type AuthService interface {
	Create(ctx context.Context, info *model.UserInfo, password, passwordConfirm string) (int64, error)
	Get(ctx context.Context, id int64) (*model.User, error)
	GetList(ctx context.Context, limit pagination.Limit) ([]model.User, error)
	Update(ctx context.Context, id int64, info *model.UpdateUserInfo) error
	Delete(ctx context.Context, id int64) error
	Login(ctx context.Context, email, password string) (*model.TokenPair, error)
	Check(ctx context.Context, endpointAddress string) (bool, error)
	GetRefreshToken(ctx context.Context, refreshToken string) (string, error)
	GetAccessToken(ctx context.Context, refreshToken string) (string, error)
}

// AuditLogService audit log service
type AuditLogService interface {
	Create(ctx context.Context, message *model.Log) error
}

// HashService hash service
type HashService interface {
	Hash(ctx context.Context, password string) (string, error)
	Check(ctx context.Context, password, hash string) bool
	GenerateSalt(ctx context.Context) (string, error)
}

// CacheUserService cache service
type CacheUserService interface {
	Set(ctx context.Context, user *model.User) error
	SetPartial(ctx context.Context, id int64, user *model.UpdateUserInfo) error
	SetList(ctx context.Context, users []model.User, limit pagination.Limit) error
	Get(ctx context.Context, userID int64) (*model.User, error)
	GetList(ctx context.Context, limit pagination.Limit) ([]model.User, error)
	Delete(ctx context.Context, userID int64) error
}

// ValidatorService validator service
type ValidatorService interface {
	ValidateUser(ctx context.Context, userInfo model.UserInfo, password, passwordConfirm string) error
}

// ConsumerService consumer service
type ConsumerService interface {
	RunConsumer(ctx context.Context) error
}

// ProducerService producer service
type ProducerService interface {
	SendMessage(ctx context.Context, message model.LogMessage) error
}

// TokenGenerator token generator
type TokenGenerator interface {
	GenerateToken(info model.UserInfo, secretKey []byte, duration time.Duration) (string, error)
}
