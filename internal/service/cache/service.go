package cache

import (
	"github.com/nqxcode/auth_microservice/internal/repository"
	"sync"

	def "github.com/nqxcode/auth_microservice/internal/service"
)

type service struct {
	mu             sync.RWMutex
	userRepository repository.UserRepository
}

// NewService new cache service
func NewService(userRepository repository.UserRepository) def.CacheService {
	return &service{
		userRepository: userRepository,
	}
}
