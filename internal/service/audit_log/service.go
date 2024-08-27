package audit_log

import (
	"github.com/nqxcode/auth_microservice/internal/repository"
	def "github.com/nqxcode/auth_microservice/internal/service"
)

type service struct {
	logRepository repository.LogRepository
}

// NewService new audit log service
func NewService(logRepository repository.LogRepository) def.AuditLogService {
	return &service{
		logRepository: logRepository,
	}
}
