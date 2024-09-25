package auth

import (
	"context"
	"github.com/opentracing/opentracing-go"

	"github.com/nqxcode/auth_microservice/internal/logger"
	desc "github.com/nqxcode/auth_microservice/pkg/auth_v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Login user
func (s *Implementation) Login(ctx context.Context, req *desc.LoginRequest) (*desc.LoginResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Login")
	defer span.Finish()

	logger.Info("Login user", zap.Any("email", req.GetEmail()))

	tokenPair, err := s.authService.Login(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cant login user: %v", err)
	}

	return &desc.LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}, nil
}
