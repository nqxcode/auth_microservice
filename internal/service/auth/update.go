package auth

import (
	"context"
	"log"

	"github.com/nqxcode/auth_microservice/internal/model"
	"github.com/nqxcode/auth_microservice/internal/service/log/constants"
)

func (s *service) Update(ctx context.Context, userID int64, info *model.UpdateUserInfo) error {
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		errTx := s.userRepository.Update(ctx, userID, info)
		if errTx != nil {
			return errTx
		}

		err := s.logService.Create(ctx, &model.Log{
			Message: constants.UserUpdated,
			Payload: struct {
				ID   int64
				Info *model.UpdateUserInfo
			}{
				ID:   userID,
				Info: info,
			},
		})

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	if info != nil {
		go func() {
			err = s.cacheService.SetPartial(ctx, userID, info)
			if err != nil {
				log.Println("cant set user to cache:", err)
			}
		}()
	}

	return nil
}
