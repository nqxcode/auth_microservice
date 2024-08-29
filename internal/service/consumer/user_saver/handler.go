package user_saver

import (
	"context"
	"encoding/json"
	"log"

	"github.com/IBM/sarama"

	"github.com/nqxcode/auth_microservice/internal/model"
)

func (s *service) UserSaveHandler(ctx context.Context, msg *sarama.ConsumerMessage) error {
	user := &model.User{}
	err := json.Unmarshal(msg.Value, user)
	if err != nil {
		return err
	}

	id, err := s.authService.Create(ctx, user)
	if err != nil {
		return err
	}

	log.Printf("User with id %d created\n", id)

	return nil
}
