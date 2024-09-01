package auth

import (
	"context"
	"github.com/pkg/errors"

	"github.com/nqxcode/auth_microservice/internal/model"
	"github.com/nqxcode/auth_microservice/internal/service/audit_log/constants"
)

// Update user
func (s *service) Update(ctx context.Context, userID int64, info *model.UpdateUserInfo) error {
	if info == nil {
		return errors.New("info is nil")
	}

	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		errTx := s.userRepository.Update(ctx, userID, info)
		if errTx != nil {
			return errTx
		}

		err := s.auditLogService.Create(ctx, &model.Log{
			Message: constants.UserUpdated,
			Payload: MakeAuditUpdatePayload(userID, info),
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
		err = s.cacheUserService.SetPartial(ctx, userID, info)
		if err != nil {
			return errors.Errorf("cant set user to cache: %v", err)
		}
		return nil
	})

	s.asyncRunner.Run(ctx, func(ctx context.Context) error {
		err = s.producerService.SendMessage(ctx, model.LogMessage{
			Message: constants.UserUpdated,
			Payload: MakeAuditUpdatePayload(userID, info),
		})
		if err != nil {
			return err
		}
		return nil
	})

	return nil
}
