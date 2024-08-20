package auth

import (
	"context"

	"github.com/nqxcode/auth_microservice/internal/model"
	"github.com/nqxcode/auth_microservice/internal/service/log/constants"
)

func (s *service) Update(ctx context.Context, id int64, info *model.UpdateUserInfo) error {
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		errTx := s.userRepository.Update(ctx, id, info)
		if errTx != nil {
			return errTx
		}

		err := s.logService.Create(ctx, &model.Log{
			Message: constants.UserUpdated,
			Payload: struct {
				ID   int64
				Info *model.UpdateUserInfo
			}{
				ID:   id,
				Info: info,
			},
		})

		if err != nil {
			return err
		}

		return nil
	})

	return err
}
