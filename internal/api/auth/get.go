package auth

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/nqxcode/auth_microservice/internal/converter"
	"github.com/nqxcode/auth_microservice/internal/logger"
	"github.com/nqxcode/auth_microservice/internal/tracing"
	desc "github.com/nqxcode/auth_microservice/pkg/auth_v1"
)

// Get user by id
func (s *Implementation) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "Get")
	if span != nil {
		defer span.Finish()
	}

	logger.Info("Get user", zap.Any("id", req.GetId()))

	user, err := s.authService.Get(ctx, req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cant get user: %v", err)
	}

	if user == nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return &desc.GetResponse{
		User: converter.ToUserFromService(user),
	}, nil
}
