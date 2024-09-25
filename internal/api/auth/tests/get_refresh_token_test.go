package tests

import (
	"context"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/nqxcode/auth_microservice/internal/api/auth"
	"github.com/nqxcode/auth_microservice/internal/service"
	serviceMocks "github.com/nqxcode/auth_microservice/internal/service/mocks"
	desc "github.com/nqxcode/auth_microservice/pkg/auth_v1"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGetRefreshToken(t *testing.T) {
	t.Parallel()

	type AuthServiceMockFunc func(mc *minimock.Controller) service.AuthService

	type input struct {
		ctx context.Context
		req *desc.GetRefreshTokenRequest
	}

	type expected struct {
		resp *desc.GetRefreshTokenResponse
		err  error
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		serviceErr = errors.New("invalid refresh token")

		req = &desc.GetRefreshTokenRequest{
			OldRefreshToken: "old-refresh-token",
		}

		resp = &desc.GetRefreshTokenResponse{
			RefreshToken: "new-refresh-token",
		}
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
				mock.GetRefreshTokenMock.Expect(ctx, "old-refresh-token").Return("new-refresh-token", nil)
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
				err: status.Errorf(codes.Internal, "cant get refresh token: %v", serviceErr),
			},
			authServiceMockFunc: func(mc *minimock.Controller) service.AuthService {
				mock := serviceMocks.NewAuthServiceMock(mc)
				mock.GetRefreshTokenMock.Expect(ctx, "old-refresh-token").Return("", serviceErr)
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

			_, err := api.GetRefreshToken(tt.input.ctx, tt.input.req)
			require.Equal(t, tt.expected.err, err)
		})
	}
}
