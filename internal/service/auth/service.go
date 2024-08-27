package auth

import (
	"github.com/nqxcode/auth_microservice/internal/repository"
	def "github.com/nqxcode/auth_microservice/internal/service"
	"github.com/nqxcode/auth_microservice/internal/service/async"
	"github.com/nqxcode/platform_common/client/db"
)

type service struct {
	userRepository   repository.UserRepository
	validatorService def.ValidatorService
	logService       def.LogService
	hashService      def.HashService
	cacheUserService def.CacheUserService
	txManager        db.TxManager
	asyncRunner      async.Runner
}

// NewService new auth service
func NewService(
	userRepository repository.UserRepository,
	validatorService def.ValidatorService,
	logService def.LogService,
	hashService def.HashService,
	cacheUserService def.CacheUserService,
	txManager db.TxManager,
	asyncRunner async.Runner,
) def.AuthService {
	return &service{
		userRepository:   userRepository,
		validatorService: validatorService,
		logService:       logService,
		hashService:      hashService,
		cacheUserService: cacheUserService,
		txManager:        txManager,
		asyncRunner:      asyncRunner,
	}
}
