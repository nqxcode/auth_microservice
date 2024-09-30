package auth

import (
	"context"
	"log"

	desc "github.com/nqxcode/auth_microservice/pkg/auth_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetAccessToken get access token
func (s *Implementation) GetAccessToken(ctx context.Context, req *desc.GetAccessTokenRequest) (*desc.GetAccessTokenResponse, error) {
	log.Printf("Get access token: %#v", req.GetRefreshToken())

	accessToken, err := s.authService.GetAccessToken(ctx, req.GetRefreshToken())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cant get access token: %v", err)
	}

	return &desc.GetAccessTokenResponse{
		AccessToken: accessToken,
	}, nil
}
