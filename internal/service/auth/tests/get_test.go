package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/nqxcode/auth_microservice/internal/model"
	"github.com/nqxcode/auth_microservice/internal/repository"
	repoMocks "github.com/nqxcode/auth_microservice/internal/repository/mocks"
	"github.com/nqxcode/auth_microservice/internal/service"
	"github.com/nqxcode/auth_microservice/internal/service/async"
	"github.com/nqxcode/auth_microservice/internal/service/auth"
	serviceSupport "github.com/nqxcode/auth_microservice/internal/service/auth/tests/support"
	"github.com/nqxcode/auth_microservice/internal/service/log/constants"
	serviceMocks "github.com/nqxcode/auth_microservice/internal/service/mocks"
	desc "github.com/nqxcode/auth_microservice/pkg/auth_v1"
	"github.com/nqxcode/platform_common/client/db"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	t.Parallel()

	type userRepositoryMock func(mc *minimock.Controller) repository.UserRepository
	type logServiceMock func(mc *minimock.Controller) service.LogService
	type hashServiceMock func(mc *minimock.Controller) service.HashService
	type cacheUserServiceMock func(mc *minimock.Controller) service.CacheUserService

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
		role         = int32(gofakeit.Number(int(desc.Role_ADMIN), int(desc.Role_USER)))
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
		name                 string
		input                input
		expected             expected
		userRepositoryMock   userRepositoryMock
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
			logServiceMock: func(mc *minimock.Controller) service.LogService {
				mock := serviceMocks.NewLogServiceMock(mc)
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
			txManagerFake:   serviceSupport.NewTxManagerFake(),
			asyncRunnerFake: serviceSupport.NewAsyncRunnerFake(),
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
			logServiceMock: func(mc *minimock.Controller) service.LogService {
				mock := serviceMocks.NewLogServiceMock(mc)
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
			txManagerFake:   serviceSupport.NewTxManagerFake(),
			asyncRunnerFake: serviceSupport.NewAsyncRunnerFake(),
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userRepoMock := tt.userRepositoryMock(mc)
			logSrvMock := tt.logServiceMock(mc)
			hashSrvMock := tt.hashServiceMock(mc)
			cacheUserSrvMock := tt.cacheUserServiceMock(mc)
			txMngFake := tt.txManagerFake
			asyncRunnerFake := tt.asyncRunnerFake

			srv := auth.NewService(userRepoMock, logSrvMock, hashSrvMock, cacheUserSrvMock, txMngFake, asyncRunnerFake)

			ar, err := srv.Get(tt.input.ctx, tt.input.userID)
			require.Equal(t, tt.expected.err, err)
			require.Equal(t, tt.expected.resp, ar)
		})
	}
}
