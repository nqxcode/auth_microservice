package cache

import (
	def "github.com/nqxcode/auth_microservice/internal/service"
	"github.com/nqxcode/platform_common/client/cache"
)

type service struct {
	cacheClient cache.Cache
}

// NewService new cache service
func NewService(cacheClient cache.Cache) def.CacheService {
	return &service{
		cacheClient: cacheClient,
	}
}
