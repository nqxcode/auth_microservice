package auth

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	desc "github.com/nqxcode/auth_microservice/pkg/auth_v1"
	"log"
)

// Delete user by id
func (s *Implementation) Delete(ctx context.Context, req *desc.DeleteRequest) (*empty.Empty, error) {
	log.Printf("Delete user: %+v", req.GetId())

	err := s.authService.Delete(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return nil, nil
}
