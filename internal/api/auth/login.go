package auth

import (
	"context"
	desc "github.com/nqxcode/auth_microservice/pkg/auth_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

// Login user
func (s *Implementation) Login(ctx context.Context, req *desc.LoginRequest) (*desc.LoginResponse, error) {
	log.Printf("Login user: %+v", req.GetEmail())

	refreshToken, err := s.authService.Login(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cant login user: %v", err)
	}

	return &desc.LoginResponse{
		RefreshToken: refreshToken,
	}, nil
}
