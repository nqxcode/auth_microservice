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

	"github.com/nqxcode/auth_microservice/internal/config"
	configMocks "github.com/nqxcode/auth_microservice/internal/config/mocks"
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

func TestGet(t *testing.T) {
	t.Parallel()

	type userRepositoryMock func(mc *minimock.Controller) repository.UserRepository
	type accessibleRoleRepositoryMock func(mc *minimock.Controller) repository.AccessibleRoleRepository
	type validatorServiceMock func(mc *minimock.Controller) service.ValidatorService
	type logServiceMock func(mc *minimock.Controller) service.AuditLogService
	type hashServiceMock func(mc *minimock.Controller) service.HashService
	type cacheUserServiceMock func(mc *minimock.Controller) service.CacheUserService
	type producerServiceMock func(mc *minimock.Controller) service.ProducerService
	type tokenGeneratorServiceMock func(mc *minimock.Controller) service.TokenGenerator
	type authConfigMock func(mc *minimock.Controller) config.AuthConfig

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

		id           = gofakeit.Int64()
		name         = gofakeit.Name()
		email        = gofakeit.Email()
		roles        = []desc.Role{desc.Role_ADMIN, desc.Role_USER}
		role         = int32(roles[rand.Int32N(int32(len(roles))-1)]) // nolint: gosec
		passwordHash = "HASH123"
		createdAt    = gofakeit.Date()

		repoErr = fmt.Errorf("repo error")
	)

	info := model.UserInfo{
		Name:  name,
		Email: email,
		Role:  role,
	}

	user := &model.User{
		ID:        id,
		Info:      info,
		Password:  passwordHash,
		CreatedAt: createdAt,
	}

	cases := []struct {
		name                         string
		input                        input
		expected                     expected
		userRepositoryMock           userRepositoryMock
		accessibleRoleRepositoryMock accessibleRoleRepositoryMock
		validatorServiceMock         validatorServiceMock
		logServiceMock               logServiceMock
		hashServiceMock              hashServiceMock
		cacheUserServiceMock         cacheUserServiceMock
		producerServiceMock          producerServiceMock
		txManagerFake                db.TxManager
		asyncRunnerFake              async.Runner
		tokenGeneratorServiceMock    tokenGeneratorServiceMock
		authConfigMock               authConfigMock
	}{
		{
			name: "success case",
			input: input{
				ctx:    ctx,
				userID: id,
			},
			expected: expected{
				err:  nil,
				resp: user,
			},
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.GetMock.Expect(ctx, id).Return(user, nil)
				return mock
			},
			accessibleRoleRepositoryMock: func(mc *minimock.Controller) repository.AccessibleRoleRepository {
				mock := repoMocks.NewAccessibleRoleRepositoryMock(mc)
				return mock
			},
			logServiceMock: func(mc *minimock.Controller) service.AuditLogService {
				mock := serviceMocks.NewAuditLogServiceMock(mc)
				mock.CreateMock.Expect(ctx, &model.Log{
					Message: constants.UserFound,
					Payload: user,
				}).Return(nil)
				return mock
			},
			hashServiceMock: func(mc *minimock.Controller) service.HashService {
				mock := serviceMocks.NewHashServiceMock(mc)
				return mock
			},
			cacheUserServiceMock: func(mc *minimock.Controller) service.CacheUserService {
				mock := serviceMocks.NewCacheUserServiceMock(mc)
				mock.GetMock.Expect(ctx, id).Return(nil, nil)
				mock.SetMock.Expect(ctx, user).Return(nil)
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
			tokenGeneratorServiceMock: func(mc *minimock.Controller) service.TokenGenerator {
				mock := serviceMocks.NewTokenGeneratorMock(mc)
				return mock
			},
			authConfigMock: func(mc *minimock.Controller) config.AuthConfig {
				return configMocks.NewAuthConfigMock(mc)
			},
		},
		{
			name: "service error case",
			input: input{
				ctx:    ctx,
				userID: id,
			},
			expected: expected{
				err:  repoErr,
				resp: (*model.User)(nil),
			},
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.GetMock.Expect(ctx, id).Return(nil, repoErr)
				return mock
			},
			accessibleRoleRepositoryMock: func(mc *minimock.Controller) repository.AccessibleRoleRepository {
				mock := repoMocks.NewAccessibleRoleRepositoryMock(mc)
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
				mock.GetMock.Expect(ctx, id).Return(nil, nil)
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
			tokenGeneratorServiceMock: func(mc *minimock.Controller) service.TokenGenerator {
				mock := serviceMocks.NewTokenGeneratorMock(mc)
				return mock
			},
			authConfigMock: func(mc *minimock.Controller) config.AuthConfig {
				return configMocks.NewAuthConfigMock(mc)
			},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userRepoMock := tt.userRepositoryMock(mc)
			accessibleRoleRepoMock := tt.accessibleRoleRepositoryMock(mc)
			validatorSrvMock := tt.validatorServiceMock(mc)
			logSrvMock := tt.logServiceMock(mc)
			hashSrvMock := tt.hashServiceMock(mc)
			cacheUserSrvMock := tt.cacheUserServiceMock(mc)
			txMngFake := tt.txManagerFake
			producerSrv := tt.producerServiceMock(mc)
			asyncRunnerFake := tt.asyncRunnerFake
			authConfig := tt.authConfigMock(mc)
			tokenGenerator := tt.tokenGeneratorServiceMock(mc)

			srv := auth.NewService(userRepoMock, accessibleRoleRepoMock, validatorSrvMock, logSrvMock, hashSrvMock, cacheUserSrvMock, txMngFake, producerSrv, asyncRunnerFake, tokenGenerator, authConfig)

			ar, err := srv.Get(tt.input.ctx, tt.input.userID)
			require.Equal(t, tt.expected.err, err)
			require.Equal(t, tt.expected.resp, ar)
		})
	}
}
