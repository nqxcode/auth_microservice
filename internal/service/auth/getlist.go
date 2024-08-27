package auth

import (
	"context"

	"github.com/nqxcode/auth_microservice/internal/model"
	"github.com/nqxcode/auth_microservice/internal/service/log/constants"

	modelCommon "github.com/nqxcode/platform_common/model"
	"github.com/nqxcode/platform_common/pagination"
	"github.com/pkg/errors"
)

// GetList get user list
func (s *service) GetList(ctx context.Context, limit pagination.Limit) ([]model.User, error) {
	var (
		users []model.User
		err   error
	)
	users, err = s.cacheUserService.GetList(ctx, limit)
	if err != nil {
		return nil, err
	}

	if len(users) > 0 {
		return users, nil
	}

	err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		users, errTx = s.userRepository.GetList(ctx, limit)
		if errTx != nil {
			return errTx
		}

		errLog := s.logService.Create(ctx, &model.Log{
			Message: constants.UserList,
			Payload: modelCommon.ExtractIDs(users),
		})

		if errLog != nil {
			return errLog
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if len(users) > 0 {
		s.asyncRunner.Run(ctx, func(ctx context.Context) error {
			err = s.cacheUserService.SetList(ctx, users, limit)
			if err != nil {
				return errors.Errorf("cant set many users to cache: %v", err)
			}
			return nil
		})
	}

	return users, nil
}
