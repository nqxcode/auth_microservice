package audit_log

import (
	"context"
	"encoding/json"

	"github.com/nqxcode/auth_microservice/internal/model"

	"github.com/nqxcode/platform_common/helper/grpc"
)

// Create audit log
func (s *service) Create(ctx context.Context, log *model.Log) error {
	ip, _ := grpc.ClientIP(ctx)
	jsonPayload, _ := json.Marshal(log.Payload)

	err := s.logRepository.Create(ctx, &model.Log{
		Message: log.Message,
		Payload: jsonPayload,
		IP:      ip,
	})
	if err != nil {
		return err
	}

	return nil
}
