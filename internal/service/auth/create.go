package auth

import (
	"context"
	"log"

	"github.com/nqxcode/auth_microservice/internal/model"
	"github.com/nqxcode/auth_microservice/internal/service/log/constants"
)

func (s *service) Create(ctx context.Context, user *model.User) (int64, error) {
	var userID int64
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		password, errHash := s.hashService.Hash(ctx, user.Password)
		if errHash != nil {
			return errHash
		}

		var errTx error
		userID, errTx = s.userRepository.Create(ctx, &model.User{Info: user.Info, Password: password})
		if errTx != nil {
			return errTx
		}

		user, errTx = s.userRepository.Get(ctx, userID)
		if errTx != nil {
			return errTx
		}

		err := s.logService.Create(ctx, &model.Log{
			Message: constants.UserCreated,
			Payload: user,
		})

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	if user != nil {
		go func() {
			err = s.cacheService.Set(ctx, user)
			if err != nil {
				log.Println("cant set user to cache:", err)
			}
		}()
	}

	return userID, nil
}
