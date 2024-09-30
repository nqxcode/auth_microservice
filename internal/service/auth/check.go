package auth

import (
	"context"
	"github.com/nqxcode/auth_microservice/internal/model"
	"strings"

	"github.com/nqxcode/auth_microservice/internal/utils"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const authPrefix = "Bearer "

// Check access to endpoint
func (s *service) Check(ctx context.Context, endpointAddress string) (bool, error) {
	// Extract authorization header
	authHeader, err := s.extractAuthHeader(ctx)
	if err != nil {
		return false, err
	}

	// Verify access token
	claims, err := s.verifyAccessToken(authHeader)
	if err != nil {
		return false, err
	}

	// check user by email
	user, err := s.getUserByEmail(ctx, claims.Email)
	if err != nil {
		return false, err
	}
	if user == nil {
		return false, errors.New("user not found")
	}

	// check user role
	hasAccess, err := s.checkUserRole(ctx, endpointAddress, claims.Role)
	if err != nil {
		return false, err
	}
	if !hasAccess {
		return false, errors.New("access denied")
	}

	return true, nil
}

// Extract authorization header
func (s *service) extractAuthHeader(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.New("metadata is not provided")
	}

	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return "", errors.New("authorization header is not provided")
	}

	if !strings.HasPrefix(authHeader[0], authPrefix) {
		return "", errors.New("invalid authorization header format")
	}

	return strings.TrimPrefix(authHeader[0], authPrefix), nil
}

// Verify access token
func (s *service) verifyAccessToken(accessToken string) (*model.UserClaims, error) {
	claims, err := utils.VerifyToken(accessToken, []byte(s.authConfig.AccessTokenSecretKey()))
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "invalid access token")
	}
	return claims, nil
}

// Get user by email
func (s *service) getUserByEmail(ctx context.Context, email string) (*model.User, error) {
	user, err := s.userRepository.GetByEmail(ctx, email)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user by email")
	}
	return user, nil
}

// Check user role
func (s *service) checkUserRole(ctx context.Context, endpointAddress, userRole string) (bool, error) {
	accessibleRoleMap, err := s.getAccessibleRoles(ctx)
	if err != nil {
		return false, errors.New("failed to get accessible roles")
	}

	roleMap, ok := accessibleRoleMap[endpointAddress]
	if !ok {
		return false, errors.New("endpoint not found")
	}

	if _, ok = roleMap[userRole]; ok {
		return true, nil
	}

	return false, nil
}

// Get accessible roles
func (s *service) getAccessibleRoles(ctx context.Context) (map[string]map[string]struct{}, error) {
	if s.accessibleRoles == nil {
		s.accessibleRoles = make(map[string]map[string]struct{})

		roles, err := s.accessibleRoleRepository.GetList(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get accessible roles")
		}

		for _, role := range roles {
			if _, ok := s.accessibleRoles[role.EndpointAddress]; !ok {
				s.accessibleRoles[role.EndpointAddress] = make(map[string]struct{})
			}
			s.accessibleRoles[role.EndpointAddress][role.Role] = struct{}{}
		}
	}

	return s.accessibleRoles, nil
}
