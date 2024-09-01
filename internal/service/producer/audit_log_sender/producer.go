package audit_log_sender

import (
	"context"
	"fmt"

	"github.com/nqxcode/platform_common/helper/grpc"

	"github.com/nqxcode/auth_microservice/internal/model"
	def "github.com/nqxcode/auth_microservice/internal/service"
	"github.com/nqxcode/platform_common/client/broker/kafka"
)

var _ def.ProducerService = (*service)(nil)

type service struct {
	syncProducer kafka.SyncProducer
}

// NewService new producer service
func NewService(syncProducer kafka.SyncProducer) def.ProducerService {
	return &service{syncProducer: syncProducer}
}

func (s service) SendMessage(ctx context.Context, message model.LogMessage) error {
	message.IP, _ = grpc.ClientIP(ctx)
	_, _, err := s.syncProducer.Produce(ctx, "logs-topic", message)
	if err != nil {
		return fmt.Errorf("syncProducer.Produce %w", err)
	}
	return nil
}
