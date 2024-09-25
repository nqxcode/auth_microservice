package auth

import (
	"context"

	"github.com/nqxcode/auth_microservice/internal/converter"
	"github.com/nqxcode/auth_microservice/internal/model"
	"github.com/nqxcode/auth_microservice/internal/utils"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *service) GetAccessToken(ctx context.Context, refreshToken string) (string, error) {
	claims, err := utils.VerifyToken(refreshToken, []byte(s.authConfig.RefreshTokenSecretKey()))
	if err != nil {
		return "", status.Errorf(codes.Aborted, "invalid refresh token")
	}

	user, err := s.userRepository.GetByEmail(ctx, claims.Email)
	if err != nil {
		return "", errors.Wrap(err, "failed to get user by email")
	}
	if user == nil {
		return "", errors.New("user not found")
	}

	accessToken, err := s.tokenGeneratorService.GenerateToken(model.UserInfo{
		Name:  claims.Username,
		Email: claims.Email,
		Role:  converter.ToRole(claims.Role),
	},
		[]byte(s.authConfig.AccessTokenSecretKey()),
		s.authConfig.AccessTokenExpiration(),
	)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}
