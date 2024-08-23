package auth

import (
	"context"
	"github.com/nqxcode/auth_microservice/internal/service/log/constants"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"log"

	"github.com/nqxcode/auth_microservice/internal/model"
	"github.com/nqxcode/platform_common/pagination"
)

func (s *service) GetList(ctx context.Context, limit *pagination.Limit) ([]model.User, error) {
	var (
		users []model.User
		err   error
	)
	users, err = s.cacheService.GetList(ctx, id)
	if err != nil {
		if !errors.Is(err, redis.Nil) { // TODO check this comparison
			return nil, err
		}
	}

	if user != nil {
		return user, nil
	}

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

	go func() {
		err = s.cacheService.SetList(ctx, users)
		if err != nil {
			log.Println("cant set many users to cache:", err)
		}
	}()

	return users, nil
}
