package auth

import (
	"context"
	"log"

	"github.com/nqxcode/auth_microservice/internal/model"
	"github.com/nqxcode/auth_microservice/internal/service/log/constants"
)

func (s *service) Delete(ctx context.Context, userID int64) error {
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		errTx := s.userRepository.Delete(ctx, userID)
		if errTx != nil {
			return errTx
		}

		err := s.logService.Create(ctx, &model.Log{
			Message: constants.UserDeleted,
			Payload: userID,
		})

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	go func() {
		err = s.cacheService.Delete(ctx, userID)
		if err != nil {
			log.Println("cant delete user to cache:", err)
		}
	}()

	return nil
}
