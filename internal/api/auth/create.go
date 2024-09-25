package auth

import (
	"context"

	"go.uber.org/zap"

	"github.com/nqxcode/auth_microservice/internal/converter"
	"github.com/nqxcode/auth_microservice/internal/logger"
	desc "github.com/nqxcode/auth_microservice/pkg/auth_v1"
)

// Create user
func (s *Implementation) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	logger.Info("Create user", zap.Any("info", req.GetInfo()))

	userID, err := s.authService.Create(ctx, converter.ToUserInfoFromDesc(req.GetInfo()), req.GetPassword(), req.GetPasswordConfirm())
	if err != nil {
		return nil, err
	}

	return &desc.CreateResponse{
		Id: userID,
	}, nil
}
