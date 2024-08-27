package tests

import (
	"context"
	"fmt"
	"math/rand/v2"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/nqxcode/platform_common/client/db"
	"github.com/stretchr/testify/require"

	"github.com/nqxcode/auth_microservice/internal/model"
	"github.com/nqxcode/auth_microservice/internal/repository"
	repoMocks "github.com/nqxcode/auth_microservice/internal/repository/mocks"
	"github.com/nqxcode/auth_microservice/internal/service"
	"github.com/nqxcode/auth_microservice/internal/service/async"
	"github.com/nqxcode/auth_microservice/internal/service/audit_log/constants"
	"github.com/nqxcode/auth_microservice/internal/service/auth"
	serviceSupport "github.com/nqxcode/auth_microservice/internal/service/auth/tests/support"
	serviceMocks "github.com/nqxcode/auth_microservice/internal/service/mocks"
	desc "github.com/nqxcode/auth_microservice/pkg/auth_v1"
)

func TestUpdate(t *testing.T) {
	t.Parallel()

	type userRepositoryMock func(mc *minimock.Controller) repository.UserRepository
	type validatorServiceMock func(mc *minimock.Controller) service.ValidatorService
	type logServiceMock func(mc *minimock.Controller) service.AuditLogService
	type hashServiceMock func(mc *minimock.Controller) service.HashService
	type cacheUserServiceMock func(mc *minimock.Controller) service.CacheUserService

	type input struct {
		ctx    context.Context
		userID int64
		info   *model.UpdateUserInfo
	}

	type expected struct {
		err error
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id    = gofakeit.Int64()
		name  = gofakeit.Name()
		roles = []desc.Role{desc.Role_ADMIN, desc.Role_USER}
		role  = int32(roles[rand.Int32N(int32(len(roles))-1)]) // nolint: gosec

		repoErr = fmt.Errorf("repo error")
	)

	info := &model.UpdateUserInfo{
		Name: &name,
		Role: &role,
	}

	cases := []struct {
		name                 string
		input                input
		expected             expected
		userRepositoryMock   userRepositoryMock
		validatorServiceMock validatorServiceMock
		logServiceMock       logServiceMock
		hashServiceMock      hashServiceMock
		cacheUserServiceMock cacheUserServiceMock
		txManagerFake        db.TxManager
		asyncRunnerFake      async.Runner
	}{
		{
			name: "success case",
			input: input{
				ctx:    ctx,
				userID: id,
				info:   info,
			},
			expected: expected{
				err: nil,
			},
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.UpdateMock.Expect(ctx, id, info).Return(nil)
				return mock
			},
			logServiceMock: func(mc *minimock.Controller) service.AuditLogService {
				mock := serviceMocks.NewAuditLogServiceMock(mc)
				mock.CreateMock.Expect(ctx, &model.Log{
					Message: constants.UserUpdated,
					Payload: struct {
						ID   int64
						Info *model.UpdateUserInfo
					}{ID: id, Info: info},
				}).Return(nil)
				return mock
			},
			hashServiceMock: func(mc *minimock.Controller) service.HashService {
				mock := serviceMocks.NewHashServiceMock(mc)
				return mock
			},
			cacheUserServiceMock: func(mc *minimock.Controller) service.CacheUserService {
				mock := serviceMocks.NewCacheUserServiceMock(mc)
				mock.SetPartialMock.Expect(ctx, id, &model.UpdateUserInfo{Name: &name, Role: &role}).Return(nil)
				return mock
			},
			txManagerFake:   serviceSupport.NewTxManagerFake(),
			asyncRunnerFake: serviceSupport.NewAsyncRunnerFake(),
			validatorServiceMock: func(mc *minimock.Controller) service.ValidatorService {
				mock := serviceMocks.NewValidatorServiceMock(mc)
				return mock
			},
		},
		{
			name: "service error case",
			input: input{
				ctx:    ctx,
				userID: id,
				info:   info,
			},
			expected: expected{
				err: repoErr,
			},
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.UpdateMock.Expect(ctx, id, info).Return(repoErr)
				return mock
			},
			logServiceMock: func(mc *minimock.Controller) service.AuditLogService {
				mock := serviceMocks.NewAuditLogServiceMock(mc)
				return mock
			},
			cacheUserServiceMock: func(mc *minimock.Controller) service.CacheUserService {
				mock := serviceMocks.NewCacheUserServiceMock(mc)
				return mock
			},
			hashServiceMock: func(mc *minimock.Controller) service.HashService {
				mock := serviceMocks.NewHashServiceMock(mc)
				return mock
			},
			txManagerFake:   serviceSupport.NewTxManagerFake(),
			asyncRunnerFake: serviceSupport.NewAsyncRunnerFake(),
			validatorServiceMock: func(mc *minimock.Controller) service.ValidatorService {
				mock := serviceMocks.NewValidatorServiceMock(mc)
				return mock
			},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userRepoMock := tt.userRepositoryMock(mc)
			validatorSrvMock := tt.validatorServiceMock(mc)
			logSrvMock := tt.logServiceMock(mc)
			hashSrvMock := tt.hashServiceMock(mc)
			cacheUserSrvMock := tt.cacheUserServiceMock(mc)
			txMngFake := tt.txManagerFake
			asyncRunnerFake := tt.asyncRunnerFake

			srv := auth.NewService(userRepoMock, validatorSrvMock, logSrvMock, hashSrvMock, cacheUserSrvMock, txMngFake, asyncRunnerFake)

			err := srv.Update(tt.input.ctx, tt.input.userID, tt.input.info)
			require.Equal(t, tt.expected.err, err)
		})
	}
}
