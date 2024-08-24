package redis

import (
	"context"
	"strconv"
	"strings"

	"github.com/nqxcode/auth_microservice/internal/model"
	"github.com/nqxcode/auth_microservice/internal/repository"
	"github.com/nqxcode/auth_microservice/internal/repository/user/redis/converter"
	modelRepo "github.com/nqxcode/auth_microservice/internal/repository/user/redis/model"
	"github.com/nqxcode/platform_common/client/cache"
	"github.com/nqxcode/platform_common/helper/slice"
	"github.com/nqxcode/platform_common/helper/time"
	"github.com/nqxcode/platform_common/pagination"

	redigo "github.com/gomodule/redigo/redis"
)

const cacheKey = "user"

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
		UpdatedAtNs: time.ToUnixNanoFromSqlNullTime(model.UpdatedAt),
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
	if info == nil {
		return nil
	}

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

	if !changes {
		return nil
	}

	err = r.redisClient.HashSet(ctx, buildCacheKey(idStr), user)
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
		return nil, nil
	}

	var user modelRepo.User
	err = redigo.ScanStruct(values, &user)
	if err != nil {
		return nil, err
	}

	return converter.ToUserFromRepo(&user), nil
}

// GetByIDs get users by ids
func (r repo) GetByIDs(ctx context.Context, ids []int64) ([]model.User, error) {
	valuesList := make([]cache.Values, 0, len(ids))

	valuesList, err := r.redisClient.MultiHGetAll(ctx, func(ids []int64) []string {
		result := make([]string, len(ids))
		for i := range ids {
			idStr := strconv.FormatInt(ids[i], 10)
			result[i] = buildCacheKey(idStr)
		}
		return result
	}(ids))
	if err != nil {
		return nil, err
	}

	users := make([]modelRepo.User, 0, len(valuesList))
	for _, v := range valuesList {
		if len(v.Values) == 0 {
			continue
		}
		var user modelRepo.User
		err = redigo.ScanStruct(v.Values, &user)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return converter.ToManyUserFromRepo(users), nil
}

// GetList users
func (r repo) GetList(ctx context.Context, limit pagination.Limit) ([]model.User, error) {
	cacheKeyPrefix := buildCacheKeyPrefix()

	keys, err := r.redisClient.Scan(ctx, buildCacheKey("*"), cache.WithKeyComparator(func(a, b string) bool {
		aNum := extractNumberAfterPrefix(a, cacheKeyPrefix)
		bNum := extractNumberAfterPrefix(b, cacheKeyPrefix)
		return aNum < bNum
	}))

	if err != nil {
		return nil, err
	}

	keys = slice.SliceByLimit(keys, limit)

	valuesList, err := r.redisClient.MultiHGetAll(ctx, keys)
	if err != nil {
		return nil, err
	}

	users := make([]modelRepo.User, 0, len(valuesList))
	for _, v := range valuesList {
		var user modelRepo.User
		err = redigo.ScanStruct(v.Values, &user)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return converter.ToManyUserFromRepo(users), nil
}

func buildCacheKey(value string) string {
	return buildCacheKeyPrefix() + value
}

func buildCacheKeyPrefix() string {
	return cacheKey + ":"
}

func extractNumberAfterPrefix(key, prefix string) int {
	num, err := strconv.Atoi(strings.TrimPrefix(key, prefix))
	if err == nil {
		return num
	}
	return 0
}
