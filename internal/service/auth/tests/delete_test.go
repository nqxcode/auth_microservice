package tests

import (
	"context"
	"fmt"
	"github.com/nqxcode/auth_microservice/internal/config"
	configMocks "github.com/nqxcode/auth_microservice/internal/config/mocks"
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
)

func TestDelete(t *testing.T) {
	t.Parallel()

	type userRepositoryMock func(mc *minimock.Controller) repository.UserRepository
	type validatorServiceMock func(mc *minimock.Controller) service.ValidatorService
	type logServiceMock func(mc *minimock.Controller) service.AuditLogService
	type hashServiceMock func(mc *minimock.Controller) service.HashService
	type cacheUserServiceMock func(mc *minimock.Controller) service.CacheUserService
	type producerServiceMock func(mc *minimock.Controller) service.ProducerService

	type input struct {
		ctx    context.Context
		userID int64
	}

	type expected struct {
		resp any
		err  error
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id = gofakeit.Int64()

		repoErr = fmt.Errorf("repo error")
	)

	cases := []struct {
		name                 string
		input                input
		expected             expected
		userRepositoryMock   userRepositoryMock
		validatorServiceMock validatorServiceMock
		logServiceMock       logServiceMock
		hashServiceMock      hashServiceMock
		cacheUserServiceMock cacheUserServiceMock
		producerServiceMock  producerServiceMock
		txManagerFake        db.TxManager
		asyncRunnerFake      async.Runner
		authConfig           config.AuthConfig
	}{
		{
			name: "success case",
			input: input{
				ctx:    ctx,
				userID: id,
			},
			expected: expected{
				err:  nil,
				resp: nil,
			},
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.DeleteMock.Expect(ctx, id).Return(nil)
				return mock
			},
			logServiceMock: func(mc *minimock.Controller) service.AuditLogService {
				mock := serviceMocks.NewAuditLogServiceMock(mc)
				mock.CreateMock.Expect(ctx, &model.Log{
					Message: constants.UserDeleted,
					Payload: id,
				}).Return(nil)
				return mock
			},
			hashServiceMock: func(mc *minimock.Controller) service.HashService {
				mock := serviceMocks.NewHashServiceMock(mc)
				return mock
			},
			cacheUserServiceMock: func(mc *minimock.Controller) service.CacheUserService {
				mock := serviceMocks.NewCacheUserServiceMock(mc)
				mock.DeleteMock.Expect(ctx, id).Return(nil)
				return mock
			},
			producerServiceMock: func(mc *minimock.Controller) service.ProducerService {
				mock := serviceMocks.NewProducerServiceMock(mc)
				mock.SendMessageMock.Expect(ctx, model.LogMessage{Message: constants.UserDeleted, Payload: id}).Return(nil)
				return mock
			},
			txManagerFake:   serviceSupport.NewTxManagerFake(),
			asyncRunnerFake: serviceSupport.NewAsyncRunnerFake(),
			validatorServiceMock: func(mc *minimock.Controller) service.ValidatorService {
				mock := serviceMocks.NewValidatorServiceMock(mc)
				return mock
			},
			authConfig: func() config.AuthConfig {
				return configMocks.NewAuthConfigMock(mc)
			}(),
		},
		{
			name: "service error case",
			input: input{
				ctx:    ctx,
				userID: id,
			},
			expected: expected{
				err:  repoErr,
				resp: int64(0),
			},
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.DeleteMock.Expect(ctx, id).Return(repoErr)
				return mock
			},
			logServiceMock: func(mc *minimock.Controller) service.AuditLogService {
				mock := serviceMocks.NewAuditLogServiceMock(mc)
				return mock
			},
			hashServiceMock: func(mc *minimock.Controller) service.HashService {
				mock := serviceMocks.NewHashServiceMock(mc)
				return mock
			},
			cacheUserServiceMock: func(mc *minimock.Controller) service.CacheUserService {
				mock := serviceMocks.NewCacheUserServiceMock(mc)
				return mock
			},
			producerServiceMock: func(mc *minimock.Controller) service.ProducerService {
				mock := serviceMocks.NewProducerServiceMock(mc)
				return mock
			},
			txManagerFake:   serviceSupport.NewTxManagerFake(),
			asyncRunnerFake: serviceSupport.NewAsyncRunnerFake(),
			validatorServiceMock: func(mc *minimock.Controller) service.ValidatorService {
				mock := serviceMocks.NewValidatorServiceMock(mc)
				return mock
			},
			authConfig: func() config.AuthConfig {
				return configMocks.NewAuthConfigMock(mc)
			}(),
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
			cacheSrvMock := tt.cacheUserServiceMock(mc)
			txMngFake := tt.txManagerFake
			producerSrv := tt.producerServiceMock(mc)
			asyncRunnerFake := tt.asyncRunnerFake
			authConfig := tt.authConfig

			srv := auth.NewService(userRepoMock, validatorSrvMock, logSrvMock, hashSrvMock, cacheSrvMock, txMngFake, producerSrv, asyncRunnerFake, authConfig)

			err := srv.Delete(tt.input.ctx, tt.input.userID)
			require.Equal(t, tt.expected.err, err)
		})
	}
}
