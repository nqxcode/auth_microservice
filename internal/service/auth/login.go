package auth

import (
	"context"

	"github.com/pkg/errors"

	"github.com/nqxcode/auth_microservice/internal/model"
	"github.com/nqxcode/auth_microservice/internal/utils"
)

func (s *service) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.userRepository.GetByEmail(ctx, email)
	if err != nil {
		return "", errors.Wrap(err, "failed to get user by email")
	}
	if user == nil {
		return "", errors.New("user not found")
	}

	checked := s.hashService.Check(ctx, password, user.Password)
	if checked == false {
		return "", errors.New("password is incorrect")
	}

	refreshToken, err := utils.GenerateToken(model.UserInfo{
		Name:  user.Info.Name,
		Email: user.Info.Email,
		Role:  user.Info.Role,
	},
		[]byte(s.authConfig.RefreshTokenSecretKey()),
		s.authConfig.RefreshTokenExpiration(),
	)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate token")
	}

	return refreshToken, nil
}
