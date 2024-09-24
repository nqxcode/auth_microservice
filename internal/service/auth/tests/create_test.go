package tests

import (
	"context"
	"fmt"
	"math/rand/v2"
	"testing"

	"github.com/nqxcode/auth_microservice/internal/config"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/nqxcode/platform_common/client/db"
	"github.com/stretchr/testify/require"

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

func TestCreate(t *testing.T) {
	t.Parallel()

	type userRepositoryMock func(mc *minimock.Controller) repository.UserRepository
	type accessibleRoleRepositoryMock func(mc *minimock.Controller) repository.AccessibleRoleRepository
	type logServiceMock func(mc *minimock.Controller) service.AuditLogService
	type hashServiceMock func(mc *minimock.Controller) service.HashService
	type cacheUserServiceMock func(mc *minimock.Controller) service.CacheUserService
	type validatorServiceMock func(mc *minimock.Controller) service.ValidatorService
	type producerServiceMock func(mc *minimock.Controller) service.ProducerService

	type input struct {
		ctx             context.Context
		info            *model.UserInfo
		password        string
		passwordConfirm string
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
		password     = gofakeit.Password(true, true, true, true, true, 8)
		passwordHash = "HASH123"

		repoErr = fmt.Errorf("repo error")
	)

	info := model.UserInfo{
		Name:  name,
		Email: email,
		Role:  role,
	}

	cases := []struct {
		name                         string
		input                        input
		expected                     expected
		userRepositoryMock           userRepositoryMock
		accessibleRoleRepositoryMock accessibleRoleRepositoryMock
		logServiceMock               logServiceMock
		hashServiceMock              hashServiceMock
		cacheUserServiceMock         cacheUserServiceMock
		validatorServiceMock         validatorServiceMock
		producerServiceMock          producerServiceMock
		txManagerFake                db.TxManager
		asyncRunnerFake              async.Runner
		authConfig                   config.AuthConfig
	}{
		{
			name: "success case",
			input: input{
				ctx:             ctx,
				info:            &info,
				password:        password,
				passwordConfirm: password,
			},
			expected: expected{
				err:  nil,
				resp: id,
			},
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.CreateMock.Expect(ctx, &model.User{Info: info, Password: passwordHash}).Return(id, nil)
				mock.GetMock.Expect(ctx, id).Return(&model.User{ID: id, Info: info, Password: passwordHash}, nil)
				return mock
			},
			accessibleRoleRepositoryMock: func(mc *minimock.Controller) repository.AccessibleRoleRepository {
				mock := repoMocks.NewAccessibleRoleRepositoryMock(mc)
				return mock
			},
			logServiceMock: func(mc *minimock.Controller) service.AuditLogService {
				mock := serviceMocks.NewAuditLogServiceMock(mc)
				mock.CreateMock.Expect(ctx, &model.Log{
					Message: constants.UserCreated,
					Payload: auth.MakeAuditCreatePayload(&model.User{ID: id, Info: info, Password: auth.HiddenPassword}),
				}).Return(nil)
				return mock
			},
			hashServiceMock: func(mc *minimock.Controller) service.HashService {
				mock := serviceMocks.NewHashServiceMock(mc)
				mock.HashMock.Expect(ctx, password).Return(passwordHash, nil)
				return mock
			},
			cacheUserServiceMock: func(mc *minimock.Controller) service.CacheUserService {
				mock := serviceMocks.NewCacheUserServiceMock(mc)
				mock.SetMock.Expect(ctx, &model.User{ID: id, Info: info, Password: passwordHash}).Return(nil)
				return mock
			},
			producerServiceMock: func(mc *minimock.Controller) service.ProducerService {
				mock := serviceMocks.NewProducerServiceMock(mc)
				mock.SendMessageMock.Expect(ctx, model.LogMessage{Message: constants.UserCreated, Payload: auth.MakeAuditCreatePayload(&model.User{ID: id, Info: info, Password: auth.HiddenPassword})}).Return(nil)
				return mock
			},
			txManagerFake:   serviceSupport.NewTxManagerFake(),
			asyncRunnerFake: serviceSupport.NewAsyncRunnerFake(),
			validatorServiceMock: func(mc *minimock.Controller) service.ValidatorService {
				mock := serviceMocks.NewValidatorServiceMock(mc)
				mock.ValidateUserMock.Expect(ctx, info, password, password).Return(nil)
				return mock
			},
			authConfig: func() config.AuthConfig {
				return configMocks.NewAuthConfigMock(mc)
			}(),
		},
		{
			name: "service error case",
			input: input{
				ctx:             ctx,
				info:            &info,
				password:        password,
				passwordConfirm: password,
			},
			expected: expected{
				err:  repoErr,
				resp: int64(0),
			},
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.CreateMock.Expect(ctx, &model.User{Info: info, Password: passwordHash}).Return(0, repoErr)
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
				mock.HashMock.Expect(ctx, password).Return(passwordHash, nil)
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
				mock.ValidateUserMock.Expect(ctx, info, password, password).Return(nil)
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
			accessibleRoleRepoMock := tt.accessibleRoleRepositoryMock(mc)
			validatorSrvMock := tt.validatorServiceMock(mc)
			logSrvMock := tt.logServiceMock(mc)
			hashSrvMock := tt.hashServiceMock(mc)
			cacheSrvMock := tt.cacheUserServiceMock(mc)
			txMngFake := tt.txManagerFake
			producerSrv := tt.producerServiceMock(mc)
			asyncRunnerFake := tt.asyncRunnerFake
			authConfig := tt.authConfig

			srv := auth.NewService(userRepoMock, accessibleRoleRepoMock, validatorSrvMock, logSrvMock, hashSrvMock, cacheSrvMock, txMngFake, producerSrv, asyncRunnerFake, authConfig)

			ar, err := srv.Create(tt.input.ctx, tt.input.info, tt.input.password, tt.input.passwordConfirm)
			require.Equal(t, tt.expected.err, err)
			require.Equal(t, tt.expected.resp, ar)
		})
	}
}
