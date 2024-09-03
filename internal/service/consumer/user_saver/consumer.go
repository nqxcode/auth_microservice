package user_saver

import (
	"context"

	def "github.com/nqxcode/auth_microservice/internal/service"
	"github.com/nqxcode/platform_common/client/broker/kafka"
)

var _ def.ConsumerService = (*service)(nil)

type service struct {
	authService def.AuthService
	consumer    kafka.Consumer
}

// NewService new consumer service
func NewService(
	authService def.AuthService,
	consumer kafka.Consumer,
) *service {
	return &service{
		authService: authService,
		consumer:    consumer,
	}
}

// RunConsumer run kafka consumer
func (s *service) RunConsumer(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-s.run(ctx):
			if err != nil {
				return err
			}
		}
	}
}

func (s *service) run(ctx context.Context) <-chan error {
	errChan := make(chan error)

	go func() {
		defer close(errChan)

		errChan <- s.consumer.Consume(ctx, UsersTopic, s.UserSaveHandler)
	}()

	return errChan
}
