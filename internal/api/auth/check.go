package auth

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/nqxcode/auth_microservice/internal/logger"
	"github.com/nqxcode/auth_microservice/internal/tracing"
	desc "github.com/nqxcode/auth_microservice/pkg/auth_v1"
)

// Check user
func (s *Implementation) Check(ctx context.Context, req *desc.CheckRequest) (*empty.Empty, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "Check")
	if span != nil {
		defer span.Finish()
	}

	logger.Info("check access token and endpoint address: %s", zap.String("endpointAddress", req.GetEndpointAddress()))

	checked, err := s.authService.Check(ctx, req.GetEndpointAddress())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cant get access: %v", err)
	}

	if !checked {
		return nil, status.Errorf(codes.PermissionDenied, "access denied")
	}

	return nil, nil
}
