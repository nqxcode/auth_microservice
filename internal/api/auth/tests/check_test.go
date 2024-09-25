package tests

import (
	"context"
	"github.com/gojuno/minimock/v3"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/nqxcode/auth_microservice/internal/api/auth"
	"github.com/nqxcode/auth_microservice/internal/service"
	serviceMocks "github.com/nqxcode/auth_microservice/internal/service/mocks"
	desc "github.com/nqxcode/auth_microservice/pkg/auth_v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func TestCheck(t *testing.T) {
	t.Parallel()

	type AuthServiceMockFunc func(mc *minimock.Controller) service.AuthService

	type input struct {
		ctx context.Context
		req *desc.CheckRequest
	}

	type expected struct {
		resp *empty.Empty
		err  error
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		req = &desc.CheckRequest{
			EndpointAddress: "/test",
		}

		resp = (*empty.Empty)(nil)
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
				mock.CheckMock.Expect(ctx, "/test").Return(true, nil)
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
				err: status.Errorf(codes.PermissionDenied, "access denied"),
			},
			authServiceMockFunc: func(mc *minimock.Controller) service.AuthService {
				mock := serviceMocks.NewAuthServiceMock(mc)
				mock.CheckMock.Expect(ctx, "/test").Return(false, nil)
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

			_, err := api.Check(tt.input.ctx, tt.input.req)
			require.Equal(t, tt.expected.err, err)
		})
	}
}
