package user

import (
	"sync"

	"github.com/nqxcode/auth_microservice/internal/repository"
	def "github.com/nqxcode/auth_microservice/internal/service"

	"github.com/nqxcode/platform_common/client/cache"
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
