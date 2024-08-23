package user

import (
	"github.com/nqxcode/platform_common/pagination"
	"strconv"
	"sync"

	"github.com/nqxcode/auth_microservice/internal/repository"
	"github.com/nqxcode/platform_common/client/cache"

	def "github.com/nqxcode/auth_microservice/internal/service"
)

const listCacheKey = "user-list"

type service struct {
	mu             sync.RWMutex
	redisClient    cache.RedisClient
	userRepository repository.UserRepository
}

// NewService new cache service
func NewService(redisClient cache.RedisClient, userRepository repository.UserRepository) def.CacheUserService {
	return &service{
		redisClient:    redisClient,
		userRepository: userRepository,
	}
}

func buildListCacheKeyByLimit(limit pagination.Limit) string {
	return buildListCacheKey(strconv.Itoa(int(limit.Offset)) + strconv.Itoa(int(limit.Limit)))
}

func buildListCacheKey(value string) string {
	return listCacheKey + ":" + value
}
