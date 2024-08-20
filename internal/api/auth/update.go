package auth

import (
	"context"
	"log"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/nqxcode/auth_microservice/internal/converter"
	desc "github.com/nqxcode/auth_microservice/pkg/auth_v1"
)

// Update user by id
func (s *Implementation) Update(ctx context.Context, req *desc.UpdateRequest) (*empty.Empty, error) {
	log.Printf("Update user info: %+v", req.GetInfo())

	err := s.authService.Update(ctx, req.GetId(), converter.ToUpdateUserInfoFromDesc(req.GetInfo()))
	if err != nil {
		return nil, err
	}

	return nil, nil
}
