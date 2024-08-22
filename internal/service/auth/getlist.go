package auth

import (
	"context"
	"github.com/nqxcode/auth_microservice/internal/service/log/constants"

	"github.com/nqxcode/auth_microservice/internal/model"
	"github.com/nqxcode/platform_common/pagination"
)

func (s *service) GetList(ctx context.Context, limit *pagination.Limit) ([]model.User, error) {
	var (
		users []model.User
		err   error
	)

	err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		users, errTx = s.userRepository.GetList(ctx, limit)
		if errTx != nil {
			return errTx
		}

		err := s.logService.Create(ctx, &model.Log{
			Message: constants.UserList,
			Payload: users,
		})

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return users, nil
}
