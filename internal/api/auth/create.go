package auth

import (
	"context"
	"log"

	"github.com/nqxcode/auth_microservice/internal/converter"
	desc "github.com/nqxcode/auth_microservice/pkg/auth_v1"
)

// Create user
func (s *Implementation) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	log.Printf("Create user: %+v", req.GetInfo())

	userID, err := s.authService.Create(ctx, converter.ToUserFromDesc(req.GetInfo(), req.GetPassword(), req.GetPasswordConfirm()))
	if err != nil {
		return nil, err
	}

	return &desc.CreateResponse{
		Id: userID,
	}, nil
}
