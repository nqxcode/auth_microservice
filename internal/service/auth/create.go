package auth

import (
	"context"

	"github.com/pkg/errors"

	"github.com/nqxcode/auth_microservice/internal/model"
	"github.com/nqxcode/auth_microservice/internal/service/log/constants"
)

// Create user
func (s *service) Create(ctx context.Context, user *model.User) (int64, error) {
	if user == nil {
		return 0, errors.New("user is nil")
	}

	if err := s.validatorService.ValidateUser(ctx, user.Info, user.Password, user.PasswordConfirm); err != nil {
		return 0, err
	}

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
		s.asyncRunner.Run(ctx, func(ctx context.Context) error {
			err = s.cacheUserService.Set(ctx, user)
			if err != nil {
				return errors.Errorf("cant set user to cache: %v", err)
			}
			return nil
		})
	}

	return userID, nil
}
