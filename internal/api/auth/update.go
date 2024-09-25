package auth

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/nqxcode/auth_microservice/internal/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/nqxcode/auth_microservice/internal/converter"
	desc "github.com/nqxcode/auth_microservice/pkg/auth_v1"
)

// Update user by id
func (s *Implementation) Update(ctx context.Context, req *desc.UpdateRequest) (*empty.Empty, error) {
	logger.Info("Update user info", zap.Any("info", req.GetInfo()))

	err := s.authService.Update(ctx, req.GetId(), converter.ToUpdateUserInfoFromDesc(req.GetInfo()))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return nil, nil
}
