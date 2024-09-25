package tests

import (
	"context"
	"fmt"
	"math/rand/v2"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"

	"github.com/nqxcode/auth_microservice/internal/config"
	configMocks "github.com/nqxcode/auth_microservice/internal/config/mocks"
	"github.com/nqxcode/auth_microservice/internal/converter"
	"github.com/nqxcode/auth_microservice/internal/model"
	"github.com/nqxcode/auth_microservice/internal/repository"
	repoMocks "github.com/nqxcode/auth_microservice/internal/repository/mocks"
	"github.com/nqxcode/auth_microservice/internal/service"
	"github.com/nqxcode/auth_microservice/internal/service/async"
	"github.com/nqxcode/auth_microservice/internal/service/auth"
	serviceSupport "github.com/nqxcode/auth_microservice/internal/service/auth/tests/support"
	serviceMocks "github.com/nqxcode/auth_microservice/internal/service/mocks"
	"github.com/nqxcode/auth_microservice/internal/utils"
	desc "github.com/nqxcode/auth_microservice/pkg/auth_v1"
	"github.com/nqxcode/platform_common/client/db"
)

func TestCheck(t *testing.T) {
	t.Parallel()

	type userRepositoryMock func(mc *minimock.Controller) repository.UserRepository
	type tokenGeneratorMock func(mc *minimock.Controller) service.TokenGenerator
	type accessibleRoleRepositoryMock func(mc *minimock.Controller) repository.AccessibleRoleRepository
	type validatorServiceMock func(mc *minimock.Controller) service.ValidatorService
	type logServiceMock func(mc *minimock.Controller) service.AuditLogService
	type hashServiceMock func(mc *minimock.Controller) service.HashService
	type cacheUserServiceMock func(mc *minimock.Controller) service.CacheUserService
	type producerServiceMock func(mc *minimock.Controller) service.ProducerService
	type tokenGeneratorServiceMock func(mc *minimock.Controller) service.TokenGenerator

	type input struct {
		ctx             context.Context
		endpointAddress string
	}

	type expected struct {
		resp bool
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

		endpointAddress = "endpoint-address"
	)

	secretKey := "secret-key" // nolint: goconst
	accessToken, err := utils.GenerateToken(model.UserInfo{
		Name:  name,
		Email: email,
		Role:  converter.ToRole(model.AdminRole),
	},
		[]byte(secretKey),
		time.Duration(1)*time.Minute,
	)
	if err != nil {
		panic(err)
	}

	incomingContext := metadata.NewIncomingContext(ctx, metadata.Pairs("authorization", "Bearer "+accessToken))

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

	accessibleRoles := []model.AccessibleRole{
		{
			EndpointAddress: endpointAddress,
			Role:            model.AdminRole,
		},
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
		tokenGeneratorMock           tokenGeneratorMock
		tokenGeneratorServiceMock    tokenGeneratorServiceMock
		authConfig                   config.AuthConfig
	}{
		{
			name: "success case",
			input: input{
				ctx:             incomingContext,
				endpointAddress: endpointAddress,
			},
			expected: expected{
				err:  nil,
				resp: true,
			},
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.GetByEmailMock.Expect(incomingContext, email).Return(user, nil)
				return mock
			},
			accessibleRoleRepositoryMock: func(mc *minimock.Controller) repository.AccessibleRoleRepository {
				mock := repoMocks.NewAccessibleRoleRepositoryMock(mc)
				mock.GetListMock.Expect(incomingContext).Return(accessibleRoles, nil)
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
				mock := configMocks.NewAuthConfigMock(mc)
				mock.AccessTokenSecretKeyMock.Return(secretKey)
				return mock
			}(),
			tokenGeneratorServiceMock: func(mc *minimock.Controller) service.TokenGenerator {
				mock := serviceMocks.NewTokenGeneratorMock(mc)
				return mock
			},
		},
		{
			name: "service error case",
			input: input{
				ctx:             incomingContext,
				endpointAddress: endpointAddress,
			},
			expected: expected{
				err:  errors.Errorf("failed to get user by email: %v", repoErr),
				resp: false,
			},
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.GetByEmailMock.Expect(incomingContext, email).Return(nil, repoErr)
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
				mock := configMocks.NewAuthConfigMock(mc)
				mock.AccessTokenSecretKeyMock.Return(secretKey)
				return mock
			}(),
			tokenGeneratorServiceMock: func(mc *minimock.Controller) service.TokenGenerator {
				mock := serviceMocks.NewTokenGeneratorMock(mc)
				return mock
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
			authConfig := tt.authConfig
			tokenGenerator := tt.tokenGeneratorServiceMock(mc)

			srv := auth.NewService(userRepoMock, accessibleRoleRepoMock, validatorSrvMock, logSrvMock, hashSrvMock, cacheUserSrvMock, txMngFake, producerSrv, asyncRunnerFake, tokenGenerator, authConfig)

			ar, checkErr := srv.Check(tt.input.ctx, tt.input.endpointAddress)
			if checkErr != nil {
				require.Equal(t, tt.expected.err.Error(), checkErr.Error())
			} else {
				require.Equal(t, nil, checkErr)
				require.Equal(t, tt.expected.resp, ar)
			}
		})
	}
}
