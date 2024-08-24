package tests

import (
	"context"
	"fmt"
	"testing"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/nqxcode/auth_microservice/internal/api/auth"
	"github.com/nqxcode/auth_microservice/internal/converter"
	"github.com/nqxcode/auth_microservice/internal/service"
	serviceMocks "github.com/nqxcode/auth_microservice/internal/service/mocks"
	desc "github.com/nqxcode/auth_microservice/pkg/auth_v1"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestUpdate(t *testing.T) {
	t.Parallel()

	type AuthServiceMockFunc func(mc *minimock.Controller) service.AuthService

	type input struct {
		ctx context.Context
		req *desc.UpdateRequest
	}

	type expected struct {
		resp any
		err  error
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id   = gofakeit.Int64()
		name = gofakeit.Name()
		role = desc.Role(gofakeit.Number(int(desc.Role_ADMIN), int(desc.Role_USER)))

		serviceErr = fmt.Errorf("service error")

		req = &desc.UpdateRequest{
			Id: id,
			Info: &desc.UpdateUserInfo{
				Name: wrapperspb.String(name),
				Role: role,
			},
		}

		resp = (*emptypb.Empty)(nil)
	)

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
				mock.UpdateMock.Expect(ctx, id, converter.ToUpdateUserInfoFromDesc(req.GetInfo())).Return(nil)
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
				err:  status.Error(codes.Internal, serviceErr.Error()),
				resp: resp,
			},
			authServiceMockFunc: func(mc *minimock.Controller) service.AuthService {
				mock := serviceMocks.NewAuthServiceMock(mc)
				mock.UpdateMock.Expect(ctx, id, converter.ToUpdateUserInfoFromDesc(req.GetInfo())).Return(serviceErr)
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

			ar, err := api.Update(tt.input.ctx, tt.input.req)
			require.Equal(t, tt.expected.err, err)
			require.Equal(t, tt.expected.resp, ar)
		})
	}
}
