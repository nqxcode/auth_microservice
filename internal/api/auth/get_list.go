package auth

import (
	"context"
	"github.com/opentracing/opentracing-go"

	"github.com/nqxcode/auth_microservice/internal/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/nqxcode/auth_microservice/internal/converter"
	desc "github.com/nqxcode/auth_microservice/pkg/auth_v1"
)

// GetList users by limit
func (s *Implementation) GetList(ctx context.Context, req *desc.GetListRequest) (*desc.GetListResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GetList")
	defer span.Finish()

	logger.Info("Get limit", zap.Any("limit", req.GetLimit()))

	users, err := s.authService.GetList(ctx, converter.ToLimitFromDesc(req.GetLimit()))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cant get users: %v", err)
	}

	return &desc.GetListResponse{
		Users: converter.ToUsersFromService(users),
	}, nil
}
