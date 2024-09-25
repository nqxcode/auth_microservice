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

// GetAccessToken get access token
func (s *Implementation) GetAccessToken(ctx context.Context, req *desc.GetAccessTokenRequest) (*desc.GetAccessTokenResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GetAccessToken")
	defer span.Finish()

	logger.Info("Get access token", zap.Any("refreshToken", req.GetRefreshToken))

	accessToken, err := s.authService.GetAccessToken(ctx, req.GetRefreshToken())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cant get access token: %v", err)
	}

	return &desc.GetAccessTokenResponse{
		AccessToken: accessToken,
	}, nil
}
