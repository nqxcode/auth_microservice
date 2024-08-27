package tests

import (
	"context"
	"fmt"
	"math/rand/v2"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	helperGob "github.com/nqxcode/platform_common/helper/gob"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/nqxcode/auth_microservice/internal/api/auth"
	"github.com/nqxcode/auth_microservice/internal/converter"
	"github.com/nqxcode/auth_microservice/internal/service"
	serviceMocks "github.com/nqxcode/auth_microservice/internal/service/mocks"
	desc "github.com/nqxcode/auth_microservice/pkg/auth_v1"
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
		roles    = []desc.Role{desc.Role_ADMIN, desc.Role_USER}
		role     = roles[rand.Int32N(int32(len(roles))-1)] // nolint: gosec
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

	invalidEmailReq := func() *desc.CreateRequest {
		var invalidReq *desc.CreateRequest
		helperGob.DeepClone(req, &invalidReq)

		invalidReq.Info.Email = "invalid email"

		return invalidReq
	}

	emptyEmailReq := func() *desc.CreateRequest {
		var invalidReq *desc.CreateRequest
		helperGob.DeepClone(req, &invalidReq)

		invalidReq.Info.Email = ""

		return invalidReq
	}

	emptyNameReq := func() *desc.CreateRequest {
		var invalidReq *desc.CreateRequest
		helperGob.DeepClone(req, &invalidReq)

		invalidReq.Info.Name = ""

		return invalidReq
	}

	diffPasswordReq := func() *desc.CreateRequest {
		var invalidReq *desc.CreateRequest
		helperGob.DeepClone(req, &invalidReq)

		invalidReq.PasswordConfirm = "321"

		return invalidReq
	}

	nilInfoReq := func() *desc.CreateRequest {
		var invalidReq *desc.CreateRequest
		helperGob.DeepClone(req, &invalidReq)

		invalidReq.Info = nil

		return invalidReq
	}

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
				mock.CreateMock.Expect(ctx, converter.ToUserFromDesc(req.GetInfo(), req.GetPassword(), req.GetPasswordConfirm())).Return(id, nil)
				return mock
			},
		},
		{
			name: "invalid input case - invalid email",
			input: input{
				ctx: ctx,
				req: invalidEmailReq(),
			},
			expected: expected{
				err: status.Error(codes.Internal, "invalid email format"),
			},
			authServiceMockFunc: func(mc *minimock.Controller) service.AuthService {
				mock := serviceMocks.NewAuthServiceMock(mc)
				mock.CreateMock.Expect(ctx, converter.ToUserFromDesc(invalidEmailReq().GetInfo(), invalidEmailReq().GetPassword(), invalidEmailReq().GetPasswordConfirm())).Return(0, status.Error(codes.Internal, "invalid email format"))
				return mock
			},
		},
		{
			name: "invalid input case - empty email",
			input: input{
				ctx: ctx,
				req: func() *desc.CreateRequest {
					var invalidReq *desc.CreateRequest
					helperGob.DeepClone(req, &invalidReq)

					invalidReq.Info.Email = ""

					return invalidReq
				}(),
			},
			expected: expected{
				err: status.Error(codes.Internal, "email is required"),
			},
			authServiceMockFunc: func(mc *minimock.Controller) service.AuthService {
				mock := serviceMocks.NewAuthServiceMock(mc)
				mock.CreateMock.Expect(ctx, converter.ToUserFromDesc(emptyEmailReq().GetInfo(), emptyEmailReq().GetPassword(), emptyEmailReq().GetPasswordConfirm())).Return(0, status.Error(codes.Internal, "email is required"))
				return mock
			},
		},
		{
			name: "invalid input case - empty name",
			input: input{
				ctx: ctx,
				req: emptyNameReq(),
			},
			expected: expected{
				err: status.Error(codes.Internal, "name is required"),
			},
			authServiceMockFunc: func(mc *minimock.Controller) service.AuthService {
				mock := serviceMocks.NewAuthServiceMock(mc)
				mock.CreateMock.Expect(ctx, converter.ToUserFromDesc(emptyNameReq().GetInfo(), emptyNameReq().GetPassword(), emptyNameReq().GetPasswordConfirm())).Return(0, status.Error(codes.Internal, "name is required"))
				return mock
			},
		},
		{
			name: "invalid input case - passwords do not match",
			input: input{
				ctx: ctx,
				req: diffPasswordReq(),
			},
			expected: expected{
				err: status.Error(codes.Internal, "passwords do not match"),
			},
			authServiceMockFunc: func(mc *minimock.Controller) service.AuthService {
				mock := serviceMocks.NewAuthServiceMock(mc)
				mock.CreateMock.Expect(ctx, converter.ToUserFromDesc(diffPasswordReq().GetInfo(), diffPasswordReq().GetPassword(), diffPasswordReq().GetPasswordConfirm())).Return(0, status.Error(codes.Internal, "passwords do not match"))

				return mock
			},
		},
		{
			name: "invalid input case - nil info",
			input: input{
				ctx: ctx,
				req: nilInfoReq(),
			},
			expected: expected{
				err: status.Error(codes.Internal, "info is required"),
			},
			authServiceMockFunc: func(mc *minimock.Controller) service.AuthService {
				mock := serviceMocks.NewAuthServiceMock(mc)
				mock.CreateMock.Expect(ctx, converter.ToUserFromDesc(nilInfoReq().GetInfo(), nilInfoReq().GetPassword(), nilInfoReq().GetPasswordConfirm())).Return(0, status.Error(codes.Internal, "info is required"))
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
				err: serviceErr,
			},
			authServiceMockFunc: func(mc *minimock.Controller) service.AuthService {
				mock := serviceMocks.NewAuthServiceMock(mc)
				mock.CreateMock.Expect(ctx, converter.ToUserFromDesc(req.GetInfo(), req.GetPassword(), req.GetPasswordConfirm())).Return(0, serviceErr)
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
