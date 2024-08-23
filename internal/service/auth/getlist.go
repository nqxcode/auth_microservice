package auth

import (
	"context"
	"log"

	"github.com/nqxcode/auth_microservice/internal/service/log/constants"
	modelCommon "github.com/nqxcode/platform_common/model"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"

	"github.com/nqxcode/auth_microservice/internal/model"
	"github.com/nqxcode/platform_common/pagination"
)

func (s *service) GetList(ctx context.Context, limit pagination.Limit) ([]model.User, error) {
	var (
		users []model.User
		err   error
	)
	users, err = s.cacheUserService.GetList(ctx, limit)
	if err != nil {
		if !errors.Is(err, redis.Nil) { // TODO check this comparison
			return nil, err
		}
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
		go func() {
			err = s.cacheUserService.SetList(ctx, users, limit)
			if err != nil {
				log.Println("cant set many users to cache:", err)
			}
		}()
	}

	return users, nil
}
