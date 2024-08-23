package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/nqxcode/auth_microservice/internal/model"
	"github.com/nqxcode/auth_microservice/internal/repository"
	repoMocks "github.com/nqxcode/auth_microservice/internal/repository/mocks"
	"github.com/nqxcode/auth_microservice/internal/service"
	"github.com/nqxcode/auth_microservice/internal/service/auth"
	testsSupport "github.com/nqxcode/auth_microservice/internal/service/auth/tests/support"
	"github.com/nqxcode/auth_microservice/internal/service/log/constants"
	serviceMocks "github.com/nqxcode/auth_microservice/internal/service/mocks"
	desc "github.com/nqxcode/auth_microservice/pkg/auth_v1"
	"github.com/nqxcode/platform_common/client/db"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	t.Parallel()
	type userRepositoryMock func(mc *minimock.Controller) repository.UserRepository
	type logServiceMock func(mc *minimock.Controller) service.LogService
	type hashServiceMock func(mc *minimock.Controller) service.HashService
	type cacheServiceMock func(mc *minimock.Controller) service.CacheUserService

	type input struct {
		ctx  context.Context
		user *model.User
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
		password     = gofakeit.Password(true, true, true, true, true, 8)
		passwordHash = "HASH123"

		repoErr = fmt.Errorf("repo error")
	)

	defer t.Cleanup(mc.Finish)

	info := model.UserInfo{
		Name:  name,
		Email: email,
		Role:  role,
	}

	cases := []struct {
		name               string
		input              input
		expected           expected
		userRepositoryMock userRepositoryMock
		logServiceMock     logServiceMock
		hashServiceMock    hashServiceMock
		cacheServiceMock   cacheServiceMock
		txManagerFake      db.TxManager
	}{
		{
			name: "success case",
			input: input{
				ctx: ctx,
				user: &model.User{
					Info:     info,
					Password: password,
				},
			},
			expected: expected{
				err:  nil,
				resp: id,
			},
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.CreateMock.Expect(ctx, &model.User{Info: info, Password: passwordHash}).Return(id, nil)
				return mock
			},
			logServiceMock: func(mc *minimock.Controller) service.LogService {
				mock := serviceMocks.NewLogServiceMock(mc)
				mock.CreateMock.Expect(ctx, &model.Log{
					Message: constants.UserCreated,
					Payload: model.User{ID: id, Info: info},
				}).Return(nil)
				return mock
			},
			hashServiceMock: func(mc *minimock.Controller) service.HashService {
				mock := serviceMocks.NewHashServiceMock(mc)
				mock.HashMock.Expect(ctx, password).Return(passwordHash, nil)
				return mock
			},
			txManagerFake: testsSupport.NewTxManagerFake(),
		},
		{
			name: "service error case",
			input: input{
				ctx: ctx,
				user: &model.User{
					Info:     info,
					Password: password,
				},
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
			logServiceMock: func(mc *minimock.Controller) service.LogService {
				mock := serviceMocks.NewLogServiceMock(mc)
				return mock
			},
			hashServiceMock: func(mc *minimock.Controller) service.HashService {
				mock := serviceMocks.NewHashServiceMock(mc)
				mock.HashMock.Expect(ctx, password).Return(passwordHash, nil)
				return mock
			},
			txManagerFake: testsSupport.NewTxManagerFake(),
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userRepoMock := tt.userRepositoryMock(mc)
			logSrvMock := tt.logServiceMock(mc)
			hashSrvMock := tt.hashServiceMock(mc)
			cacheSrvMock := tt.cacheServiceMock(mc)
			txMngFake := tt.txManagerFake

			srv := auth.NewService(userRepoMock, logSrvMock, hashSrvMock, cacheSrvMock, txMngFake)

			ar, err := srv.Create(tt.input.ctx, tt.input.user)
			require.Equal(t, tt.expected.err, err)
			require.Equal(t, tt.expected.resp, ar)
		})
	}
}
