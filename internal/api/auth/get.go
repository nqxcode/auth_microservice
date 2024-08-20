package auth

import (
	"context"
	"github.com/nqxcode/auth_microservice/internal/converter"
	desc "github.com/nqxcode/auth_microservice/pkg/auth_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

// Get user by id
func (s *Implementation) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	log.Printf("Get user: %d", req.GetId())

	user, err := s.authService.Find(ctx, req.GetId())
	if err != nil {
		return nil, status.Error(codes.NotFound, "User not found")
	}

	return &desc.GetResponse{
		User: converter.ToUserFromService(user),
	}, nil
}
