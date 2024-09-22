package auth

import (
	"github.com/nqxcode/auth_microservice/internal/config"
	"github.com/nqxcode/auth_microservice/internal/repository"
	def "github.com/nqxcode/auth_microservice/internal/service"
	"github.com/nqxcode/auth_microservice/internal/service/async"

	"github.com/nqxcode/platform_common/client/db"
)

type service struct {
	userRepository           repository.UserRepository
	accessibleRoleRepository repository.AccessibleRoleRepository
	validatorService         def.ValidatorService
	auditLogService          def.AuditLogService
	hashService              def.HashService
	cacheUserService         def.CacheUserService
	txManager                db.TxManager
	producerService          def.ProducerService
	asyncRunner              async.Runner
	authConfig               config.AuthConfig
	accessibleRoles          map[string]string
}

// NewService new auth service
func NewService(
	userRepository repository.UserRepository,
	accessibleRoleRepository repository.AccessibleRoleRepository,
	validatorService def.ValidatorService,
	auditLogService def.AuditLogService,
	hashService def.HashService,
	cacheUserService def.CacheUserService,
	txManager db.TxManager,
	producerService def.ProducerService,
	asyncRunner async.Runner,
	authConfig config.AuthConfig,
) def.AuthService {
	return &service{
		userRepository:           userRepository,
		accessibleRoleRepository: accessibleRoleRepository,
		validatorService:         validatorService,
		auditLogService:          auditLogService,
		hashService:              hashService,
		cacheUserService:         cacheUserService,
		txManager:                txManager,
		producerService:          producerService,
		asyncRunner:              asyncRunner,
		authConfig:               authConfig,
		accessibleRoles:          make(map[string]string),
	}
}
