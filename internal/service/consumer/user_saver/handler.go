package user_saver

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/nqxcode/auth_microservice/internal/logger"
	"go.uber.org/zap"

	"github.com/nqxcode/auth_microservice/internal/converter"
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
		logger.Info("Error creating user", zap.Error(err))
	} else {
		logger.Info(fmt.Sprintf("User with id %d created", id), zap.Error(err))
	}

	return nil
}
