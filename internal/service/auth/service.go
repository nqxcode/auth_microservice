package auth

import (
	"github.com/nqxcode/auth_microservice/internal/client/db"
	"github.com/nqxcode/auth_microservice/internal/repository"
	def "github.com/nqxcode/auth_microservice/internal/service"
)

type service struct {
	userRepository repository.UserRepository
	logService     def.LogService
	hashService    def.HashService
	txManager      db.TxManager
}

// NewService new auth service
func NewService(
	userRepository repository.UserRepository,
	logService def.LogService,
	hashService def.HashService,
	txManager db.TxManager,
) def.AuthService {
	return &service{
		userRepository: userRepository,
		logService:     logService,
		hashService:    hashService,
		txManager:      txManager,
	}
}
