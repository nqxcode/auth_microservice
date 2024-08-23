package auth

import (
	"context"
	"errors"
	"log"

	"github.com/nqxcode/auth_microservice/internal/model"
	"github.com/nqxcode/auth_microservice/internal/service/log/constants"

	"github.com/redis/go-redis/v9"
)

func (s *service) Get(ctx context.Context, id int64) (*model.User, error) {
	var (
		user *model.User
		err  error
	)

	user, err = s.cacheService.Get(ctx, id)
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
		user, errTx = s.userRepository.Get(ctx, id)
		if errTx != nil {
			return errTx
		}

		err := s.logService.Create(ctx, &model.Log{
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
		go func() {
			err = s.cacheService.Set(ctx, user)
			if err != nil {
				log.Println("cant set user to cache:", err)
			}
		}()
	}

	return user, err
}
