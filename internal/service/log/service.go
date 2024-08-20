package log

import (
	"github.com/nqxcode/auth_microservice/internal/repository"
	def "github.com/nqxcode/auth_microservice/internal/service"
)

type service struct {
	logRepository repository.LogRepository
}

func NewService(logRepository repository.LogRepository) def.LogService {
	return &service{
		logRepository: logRepository,
	}
}
