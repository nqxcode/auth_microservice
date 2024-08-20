package auth

import (
	"context"
	"log"
	"regexp"

	"github.com/nqxcode/auth_microservice/internal/converter"
	desc "github.com/nqxcode/auth_microservice/pkg/auth_v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Create user
func (s *Implementation) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	log.Printf("Create user: %+v", req.GetInfo())

	if err := createRequestValidate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	userID, err := s.authService.Create(ctx, converter.ToUserFromDesc(req.GetInfo(), req.GetPassword()))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &desc.CreateResponse{
		Id: userID,
	}, nil
}

// createRequestValidate validates the create request
func createRequestValidate(req *desc.CreateRequest) error {
	if req.Info == nil {
		return status.Error(codes.InvalidArgument, "info is required")
	}

	if req.Info.Name == "" {
		return status.Error(codes.InvalidArgument, "name is required")
	}

	if req.Info.Email == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	if !validateEmail(req.Info.Email) {
		return status.Error(codes.InvalidArgument, "invalid email format")
	}

	if req.Info.Role == 0 {
		return status.Error(codes.InvalidArgument, "role is required")
	}

	if req.Password != req.PasswordConfirm {
		return status.Error(codes.InvalidArgument, "passwords do not match")
	}

	return nil
}

// validateEmail checks if the provided email address is valid
func validateEmail(email string) bool {
	return regexp.
		MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`).
		MatchString(email)
}
