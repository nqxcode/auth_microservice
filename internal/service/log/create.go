package log

import (
	"context"
	"encoding/json"
	"github.com/nqxcode/auth_microservice/internal/helper"
	"github.com/nqxcode/auth_microservice/internal/model"
)

func (s *service) Create(ctx context.Context, log *model.Log) error {
	ip, _ := helper.ClientIP(ctx)
	jsonPayload, _ := json.Marshal(log.Payload)

	return s.logRepository.Create(ctx, &model.Log{
		Message: log.Message,
		Payload: jsonPayload,
		Ip:      ip,
	})
}
