package auth

import (
	"context"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	desc "github.com/nqxcode/auth_microservice/pkg/auth_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
)

// Get user by id
func (s *Implementation) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	log.Printf("User id: %d", req.GetId())

	builderSelect := sq.Select("user_id", "name", "email", "role", "created_at", "updated_at").
		From("\"user\"").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"user_id": req.GetId()}).
		Limit(1)

	query, args, err := builderSelect.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	var user model.User
	err = s.pool.QueryRow(ctx, query, args...).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		log.Fatalf("failed to select users: %v", err)
	}

	var outUpdatedAt *timestamppb.Timestamp
	if user.UpdatedAt.Valid {
		outUpdatedAt = timestamppb.New(user.UpdatedAt.Time)
	}

	return &desc.GetResponse{
		User: &desc.User{
			Id: user.ID,
			Info: &desc.UserInfo{
				Name:  user.Name,
				Email: user.Email,
				Role:  desc.Role(user.Role),
			},
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: outUpdatedAt,
		},
	}, nil
}
