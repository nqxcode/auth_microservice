package auth

import (
	"context"

	"github.com/pkg/errors"

	"github.com/nqxcode/auth_microservice/internal/model"
	"github.com/nqxcode/auth_microservice/internal/service/audit_log/constants"
)

// Delete user by id
func (s *service) Delete(ctx context.Context, userID int64) error {
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		errTx := s.userRepository.Delete(ctx, userID)
		if errTx != nil {
			return errTx
		}

		err := s.logService.Create(ctx, &model.Log{
			Message: constants.UserDeleted,
			Payload: userID,
		})

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	s.asyncRunner.Run(ctx, func(ctx context.Context) error {
		err = s.cacheUserService.Delete(ctx, userID)
		if err != nil {
			return errors.Errorf("cant delete user to cache: %v", err)
		}
		return nil
	})

	return nil
}
