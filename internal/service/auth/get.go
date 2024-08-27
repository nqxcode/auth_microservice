package auth

import (
	"context"

	"github.com/pkg/errors"

	"github.com/nqxcode/auth_microservice/internal/model"
	"github.com/nqxcode/auth_microservice/internal/service/log/constants"
)

// Get user by id
func (s *service) Get(ctx context.Context, id int64) (*model.User, error) {
	var (
		user *model.User
		err  error
	)

	user, err = s.cacheUserService.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if user != nil {
		return user, nil
	}

	err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		user, errTx = s.userRepository.Get(ctx, id)
		if errTx != nil {
			return errTx
		}

		err = s.logService.Create(ctx, &model.Log{
			Message: constants.UserFound,
			Payload: user,
		})

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
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

	return user, err
}
