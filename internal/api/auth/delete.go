package auth

import (
	"context"

	"github.com/nqxcode/auth_microservice/internal/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/golang/protobuf/ptypes/empty"
	desc "github.com/nqxcode/auth_microservice/pkg/auth_v1"
)

// Delete user by id
func (s *Implementation) Delete(ctx context.Context, req *desc.DeleteRequest) (*empty.Empty, error) {
	logger.Info("Delete user", zap.Any("id", req.GetId()))

	err := s.authService.Delete(ctx, req.GetId())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return nil, nil
}
