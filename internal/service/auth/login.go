package auth

import (
	"context"
	"github.com/pkg/errors"

	"github.com/nqxcode/auth_microservice/internal/model"
)

// Login user
func (s *service) Login(ctx context.Context, email, password string) (*model.TokenPair, error) {
	user, err := s.userRepository.GetByEmail(ctx, email)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user by email")
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	checked := s.hashService.Check(ctx, password, user.Password)
	if !checked {
		return nil, errors.New("password is incorrect")
	}

	refreshToken, err := s.tokenGeneratorService.GenerateToken(model.UserInfo{
		Name:  user.Info.Name,
		Email: user.Info.Email,
		Role:  user.Info.Role,
	},
		[]byte(s.authConfig.RefreshTokenSecretKey()),
		s.authConfig.RefreshTokenExpiration(),
	)
	if err != nil {
		return nil, err
	}

	accessToken, err := s.tokenGeneratorService.GenerateToken(model.UserInfo{
		Name:  user.Info.Name,
		Email: user.Info.Email,
		Role:  user.Info.Role,
	},
		[]byte(s.authConfig.AccessTokenSecretKey()),
		s.authConfig.AccessTokenExpiration(),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate access token")
	}

	return &model.TokenPair{RefreshToken: refreshToken, AccessToken: accessToken}, nil
}
