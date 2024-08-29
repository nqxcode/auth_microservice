package user_saver

import (
	"context"
	"encoding/json"
	"github.com/nqxcode/auth_microservice/internal/converter"
	"log"

	"github.com/IBM/sarama"

	"github.com/nqxcode/auth_microservice/internal/model"
)

func (s *service) UserSaveHandler(ctx context.Context, msg *sarama.ConsumerMessage) error {
	userMessage := &model.UserMessage{}
	err := json.Unmarshal(msg.Value, userMessage)
	if err != nil {
		return err
	}

	id, err := s.authService.Create(ctx, converter.ToUserInfoFromMessage(&userMessage.Info), userMessage.Password, userMessage.PasswordConfirm)
	if err != nil {
		log.Printf("Error creating user: %v", err)
	} else {
		log.Printf("User with id %d created\n", id)
	}

	return nil
}
