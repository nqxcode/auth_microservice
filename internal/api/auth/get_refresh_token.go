package auth

import (
	"context"

	"github.com/nqxcode/auth_microservice/internal/logger"
	desc "github.com/nqxcode/auth_microservice/pkg/auth_v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetRefreshToken get refresh token
func (s *Implementation) GetRefreshToken(ctx context.Context, req *desc.GetRefreshTokenRequest) (*desc.GetRefreshTokenResponse, error) {
	logger.Info("Get refresh token", zap.Any("oldRefreshToken", req.GetOldRefreshToken()))

	refreshToken, err := s.authService.GetRefreshToken(ctx, req.GetOldRefreshToken())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cant get refresh token: %v", err)
	}

	return &desc.GetRefreshTokenResponse{
		RefreshToken: refreshToken,
	}, nil
}
