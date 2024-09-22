package auth

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	desc "github.com/nqxcode/auth_microservice/pkg/auth_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

// Check user
func (s *Implementation) Check(ctx context.Context, req *desc.CheckRequest) (*empty.Empty, error) {
	log.Printf("Get access token: %+v", req.GetAccessToken())

	checked, err := s.authService.Check(ctx, req.GetAccessToken())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cant get access: %v", err)
	}

	if !checked {
		return nil, status.Errorf(codes.PermissionDenied, "access denied")
	}

	return nil, nil
}
