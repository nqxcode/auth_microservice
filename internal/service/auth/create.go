package auth

import (
	"context"
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

		err := s.logService.Create(ctx, &model.Log{
			Message: constants.UserCreated,
			Payload: model.User{ID: userID, Info: user.Info},
		})

		if err != nil {
			return err
		}

		return nil
	})

	return userID, err
}
