package auth

import (
	"context"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/nqxcode/auth_microservice/internal/converter"
	desc "github.com/nqxcode/auth_microservice/pkg/auth_v1"
)

// GetList users by limit
func (s *Implementation) GetList(ctx context.Context, req *desc.GetListRequest) (*desc.GetListResponse, error) {
	log.Printf("Get limit: %#v", req.GetLimit())

	users, err := s.authService.GetList(ctx, converter.ToLimitFromDesc(req.GetLimit()))
	if err != nil {
		return nil, status.Error(codes.Internal, "cant get users")
	}

	return &desc.GetListResponse{
		Users: converter.ToUsersFromService(users),
	}, nil
}
