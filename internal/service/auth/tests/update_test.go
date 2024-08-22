package tests

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/nqxcode/auth_microservice/internal/model"
	"github.com/nqxcode/auth_microservice/internal/repository"
	repoMocks "github.com/nqxcode/auth_microservice/internal/repository/mocks"
	"github.com/nqxcode/auth_microservice/internal/service"
	"github.com/nqxcode/auth_microservice/internal/service/auth"
	testsService "github.com/nqxcode/auth_microservice/internal/service/auth/tests/service"
	"github.com/nqxcode/auth_microservice/internal/service/log/constants"
	serviceMocks "github.com/nqxcode/auth_microservice/internal/service/mocks"
	desc "github.com/nqxcode/auth_microservice/pkg/auth_v1"
	"github.com/nqxcode/platform_common/client/db"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUpdate(t *testing.T) {
	t.Parallel()
	type userRepositoryMock func(mc *minimock.Controller) repository.UserRepository
	type logServiceMock func(mc *minimock.Controller) service.LogService
	type hashServiceMock func(mc *minimock.Controller) service.HashService

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

		id   = gofakeit.Int64()
		name = gofakeit.Name()
		role = int32(gofakeit.Number(int(desc.Role_ADMIN), int(desc.Role_USER)))

		repoErr = fmt.Errorf("repo error")
	)

	defer t.Cleanup(mc.Finish)

	info := &model.UpdateUserInfo{
		Name: &name,
		Role: &role,
	}

	cases := []struct {
		name               string
		input              input
		expected           expected
		userRepositoryMock userRepositoryMock
		logServiceMock     logServiceMock
		hashServiceMock    hashServiceMock
		txManagerFake      db.TxManager
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
			logServiceMock: func(mc *minimock.Controller) service.LogService {
				mock := serviceMocks.NewLogServiceMock(mc)
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
			txManagerFake: testsService.NewTxManagerFake(),
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
			logServiceMock: func(mc *minimock.Controller) service.LogService {
				mock := serviceMocks.NewLogServiceMock(mc)
				return mock
			},
			hashServiceMock: func(mc *minimock.Controller) service.HashService {
				mock := serviceMocks.NewHashServiceMock(mc)
				return mock
			},
			txManagerFake: testsService.NewTxManagerFake(),
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userRepoMock := tt.userRepositoryMock(mc)
			logSrvMock := tt.logServiceMock(mc)
			hashSrvMock := tt.hashServiceMock(mc)
			txMngFake := tt.txManagerFake

			srv := auth.NewService(userRepoMock, logSrvMock, hashSrvMock, txMngFake)

			err := srv.Update(tt.input.ctx, tt.input.userID, tt.input.info)
			require.Equal(t, tt.expected.err, err)
		})
	}
}
