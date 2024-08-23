package auth

import (
	"github.com/nqxcode/auth_microservice/internal/repository"
	def "github.com/nqxcode/auth_microservice/internal/service"
	"github.com/nqxcode/platform_common/client/db"
)

type service struct {
	userRepository   repository.UserRepository
	logService       def.LogService
	hashService      def.HashService
	cacheUserService def.CacheUserService
	txManager        db.TxManager
}

// NewService new auth service
func NewService(
	userRepository repository.UserRepository,
	logService def.LogService,
	hashService def.HashService,
	cacheUserService def.CacheUserService,
	txManager db.TxManager,
) def.AuthService {
	return &service{
		userRepository:   userRepository,
		logService:       logService,
		hashService:      hashService,
		cacheUserService: cacheUserService,
		txManager:        txManager,
	}
}
