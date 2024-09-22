package auth

import (
	"context"
	"github.com/nqxcode/auth_microservice/internal/utils"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strings"
)

const authPrefix = "Bearer "

func (s *service) Check(ctx context.Context, endpointAddress string) (bool, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return false, errors.New("metadata is not provided")
	}

	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return false, errors.New("authorization header is not provided")
	}

	if !strings.HasPrefix(authHeader[0], authPrefix) {
		return false, errors.New("invalid authorization header format")
	}

	accessToken := strings.TrimPrefix(authHeader[0], authPrefix)

	claims, err := utils.VerifyToken(accessToken, []byte(s.authConfig.AccessTokenSecretKey()))
	if err != nil {
		return false, status.Errorf(codes.Aborted, "invalid access token")
	}

	user, err := s.userRepository.GetByEmail(ctx, claims.Email)
	if err != nil {
		return false, errors.Wrap(err, "failed to get user by email")
	}
	if user == nil {
		return false, errors.New("user not found")
	}

	accessibleMap, err := s.getAccessibleRoles(ctx)
	if err != nil {
		return false, errors.New("failed to get accessible roles")
	}

	role, ok := accessibleMap[endpointAddress]
	if !ok {
		return false, nil
	}

	if role == claims.Role {
		return false, nil
	}

	return false, errors.New("access denied")
}

func (s *service) getAccessibleRoles(ctx context.Context) (map[string]string, error) {
	if s.accessibleRoles == nil {
		s.accessibleRoles = make(map[string]string)

		roles, err := s.accessibleRoleRepository.GetList(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get accessible roles")
		}

		for _, role := range roles {
			s.accessibleRoles[role.EndpointAddress] = role.Role
		}
	}

	return s.accessibleRoles, nil
}
