package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/nqxcode/auth_microservice/internal/api/auth"
	"github.com/nqxcode/auth_microservice/internal/converter"
	"github.com/nqxcode/auth_microservice/internal/service"
	serviceMocks "github.com/nqxcode/auth_microservice/internal/service/mocks"
	desc "github.com/nqxcode/auth_microservice/pkg/auth_v1"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	type AuthServiceMockFunc func(mc *minimock.Controller) service.AuthService

	type input struct {
		ctx context.Context
		req *desc.CreateRequest
	}

	type expected struct {
		resp *desc.CreateResponse
		err  error
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id       = gofakeit.Int64()
		name     = gofakeit.Name()
		email    = gofakeit.Email()
		role     = desc.Role(gofakeit.Number(int(desc.Role_ADMIN), int(desc.Role_ADMIN)))
		password = gofakeit.Password(true, true, true, true, true, 8)

		serviceErr = fmt.Errorf("service error")

		req = &desc.CreateRequest{
			Info: &desc.UserInfo{
				Name:  name,
				Email: email,
				Role:  role,
			},
			Password:        password,
			PasswordConfirm: password,
		}

		resp = &desc.CreateResponse{
			Id: id,
		}
	)

	defer t.Cleanup(mc.Finish)

	cases := []struct {
		name                string
		input               input
		expected            expected
		authServiceMockFunc AuthServiceMockFunc
	}{
		{
			name: "success case",
			input: input{
				ctx: ctx,
				req: req,
			},
			expected: expected{
				resp: resp,
			},
			authServiceMockFunc: func(mc *minimock.Controller) service.AuthService {
				mock := serviceMocks.NewAuthServiceMock(mc)
				mock.CreateMock.Expect(ctx, converter.ToUserFromDesc(req.GetInfo(), req.GetPassword())).Return(id, nil)
				return mock
			},
		},
		{
			name: "service error case",
			input: input{
				ctx: ctx,
				req: req,
			},
			expected: expected{
				err: status.Error(codes.Internal, serviceErr.Error()),
			},
			authServiceMockFunc: func(mc *minimock.Controller) service.AuthService {
				mock := serviceMocks.NewAuthServiceMock(mc)
				mock.CreateMock.Expect(ctx, converter.ToUserFromDesc(req.GetInfo(), req.GetPassword())).Return(0, serviceErr)
				return mock
			},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			authServiceMock := tt.authServiceMockFunc(mc)
			api := auth.NewImplementation(authServiceMock)

			ar, err := api.Create(tt.input.ctx, tt.input.req)
			require.Equal(t, tt.expected.err, err)
			require.Equal(t, tt.expected.resp, ar)
		})
	}
}
