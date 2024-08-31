package auth

import (
	"context"
	"github.com/nqxcode/platform_common/helper/gob"
	"github.com/pkg/errors"

	"github.com/nqxcode/auth_microservice/internal/converter"
	"github.com/nqxcode/auth_microservice/internal/model"
	"github.com/nqxcode/auth_microservice/internal/service/audit_log/constants"
)

const (
	hiddenPassword = "***"
)

// Create user
func (s *service) Create(ctx context.Context, info *model.UserInfo, password, passwordConfirm string) (int64, error) {
	if info == nil {
		return 0, errors.New("user info is nil")
	}

	if err := s.validatorService.ValidateUser(ctx, *info, password, passwordConfirm); err != nil {
		return 0, err
	}

	var (
		user   *model.User
		userID int64
	)

	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		passwordHash, errHash := s.hashService.Hash(ctx, password)
		if errHash != nil {
			return errHash
		}

		var errTx error
		userID, errTx = s.userRepository.Create(ctx, &model.User{Info: *info, Password: passwordHash})
		if errTx != nil {
			return errTx
		}

		user, errTx = s.userRepository.Get(ctx, userID)
		if errTx != nil {
			return errTx
		}

		err := s.auditLogService.Create(ctx, &model.Log{
			Message: constants.UserCreated,
			Payload: makeAuditCreatePayload(user),
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

		s.asyncRunner.Run(ctx, func(ctx context.Context) error {
			err = s.producerService.SendMessage(ctx, model.LogMessage{
				Message: constants.UserCreated,
				Payload: makeAuditCreatePayload(user),
			})
			if err != nil {
				return err
			}
			return nil
		})
	}

	return userID, nil
}

func makeAuditCreatePayload(user *model.User) any {
	if user == nil {
		return nil
	}

	var u *model.User
	gob.DeepClone(user, &u)
	u.Password = hiddenPassword

	return converter.ToLogUserMessageFromService(u)
}
