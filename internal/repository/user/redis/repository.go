package redis

import (
	"context"
	"github.com/nqxcode/auth_microservice/internal/model"
	"github.com/nqxcode/auth_microservice/internal/repository"
	"github.com/nqxcode/auth_microservice/internal/repository/user/redis/converter"
	modelRepo "github.com/nqxcode/auth_microservice/internal/repository/user/redis/model"
	"github.com/nqxcode/platform_common/client/cache"
	"github.com/nqxcode/platform_common/pagination"
	"strconv"

	redigo "github.com/gomodule/redigo/redis"
)

type repo struct {
	redisClient cache.RedisClient
}

// NewRepository new user repository
func NewRepository(redisClient cache.RedisClient) repository.UserRepository {
	return &repo{redisClient: redisClient}
}

// Create user
func (r repo) Create(ctx context.Context, model *model.User) (int64, error) {
	id := model.ID
	user := modelRepo.User{
		ID:          id,
		Name:        model.Info.Name,
		Email:       model.Info.Email,
		Role:        model.Info.Role,
		Password:    model.Password,
		CreatedAtNs: model.CreatedAt.UnixNano(),
		UpdatedAtNs: func() *int64 {
			if !model.UpdatedAt.Valid {
				return nil
			}
			return toPtr(model.UpdatedAt.Time.UnixNano())
		}(),
	}

	idStr := strconv.FormatInt(id, 10)
	err := r.redisClient.HashSet(ctx, buildCacheKey(idStr), user)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// Update user
func (r repo) Update(ctx context.Context, id int64, info *model.UpdateUserInfo) error {
	idStr := strconv.FormatInt(id, 10)
	values, err := r.redisClient.HGetAll(ctx, buildCacheKey(idStr))
	if err != nil {
		return err
	}

	if len(values) == 0 {
		return model.ErrorNoteNotFound
	}

	var user modelRepo.User
	err = redigo.ScanStruct(values, &user)
	if err != nil {
		return err
	}

	var changes bool
	if info != nil {
		if info.Role != nil {
			if user.Role != *info.Role {
				changes = true
			}
			user.Role = *info.Role
		}
		if info.Name != nil {
			if user.Name != *info.Name {
				changes = true
			}
			user.Name = *info.Name
		}
	}

	if !changes {
		return nil
	}

	err = r.redisClient.HashSet(ctx, buildCacheKey(idStr), user) // TODO set only changes field without getting of source values
	if err != nil {
		return err
	}

	return nil
}

// Delete user
func (r repo) Delete(ctx context.Context, id int64) error {
	idStr := strconv.FormatInt(id, 10)
	err := r.redisClient.Delete(ctx, buildCacheKey(idStr))
	if err != nil {
		return err
	}

	return nil
}

// Get user
func (r repo) Get(ctx context.Context, id int64) (*model.User, error) {
	idStr := strconv.FormatInt(id, 10)
	values, err := r.redisClient.HGetAll(ctx, buildCacheKey(idStr))
	if err != nil {
		return nil, err
	}

	if len(values) == 0 {
		return nil, model.ErrorNoteNotFound
	}

	var user modelRepo.User
	err = redigo.ScanStruct(values, &user)
	if err != nil {
		return nil, err
	}

	return converter.ToUserFromRepo(&user), nil
}

// GetList users
func (r repo) GetList(ctx context.Context, limit *pagination.Limit) ([]model.User, error) {
	keys, err := r.redisClient.Scan(ctx, buildCacheKey("*"))
	if err != nil {
		return nil, err
	}

	offset := limit.Offset
	end := limit.Offset + limit.Limit

	keys = keys[offset:end]

	valuesList := make([]interface{}, 0, len(keys))
	for _, key := range keys {
		values, err := r.redisClient.HGetAll(ctx, key) // TODO make this parallel
		if err != nil {
			return nil, err
		}

		if len(values) == 0 {
			return nil, model.ErrorNoteNotFound
		}

		valuesList = append(valuesList, values)
	}

	users := make([]modelRepo.User, 0, len(valuesList))
	err = redigo.ScanSlice(valuesList, &users)
	if err != nil {
		return nil, err
	}

	return converter.ToManyUserFromRepo(users), nil
}

func toPtr[T any](s T) *T {
	return &s
}

func buildCacheKey(value string) string {
	return "user:" + value
}
