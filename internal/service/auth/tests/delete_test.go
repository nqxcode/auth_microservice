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
	service2 "github.com/nqxcode/auth_microservice/internal/service/auth/tests/support"
	"github.com/nqxcode/auth_microservice/internal/service/log/constants"
	serviceMocks "github.com/nqxcode/auth_microservice/internal/service/mocks"
	"github.com/nqxcode/platform_common/client/db"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
)

func TestDelete(t *testing.T) {
	t.Parallel()
	type userRepositoryMock func(mc *minimock.Controller) repository.UserRepository
	type logServiceMock func(mc *minimock.Controller) service.LogService
	type hashServiceMock func(mc *minimock.Controller) service.HashService
	type cacheServiceMock func(mc *minimock.Controller) service.CacheService

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

	defer t.Cleanup(mc.Finish)

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
			logServiceMock: func(mc *minimock.Controller) service.LogService {
				mock := serviceMocks.NewLogServiceMock(mc)
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
			txManagerFake: service2.NewTxManagerFake(),
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
			logServiceMock: func(mc *minimock.Controller) service.LogService {
				mock := serviceMocks.NewLogServiceMock(mc)
				return mock
			},
			hashServiceMock: func(mc *minimock.Controller) service.HashService {
				mock := serviceMocks.NewHashServiceMock(mc)
				return mock
			},
			txManagerFake: service2.NewTxManagerFake(),
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

			err := srv.Delete(tt.input.ctx, tt.input.userID)
			require.Equal(t, tt.expected.err, err)
		})
	}
}
