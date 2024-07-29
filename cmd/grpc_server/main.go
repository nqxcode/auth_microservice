package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"log"
	"net"
	"time"

	desc "github.com/nqxcode/auth_microservice/pkg/auth_v1"
	"github.com/nqxcode/auth_microservice/pkg/hashing"

	sq "github.com/Masterminds/squirrel"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nqxcode/auth_microservice/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

type server struct {
	desc.UnimplementedAuthV1Server
	pool *pgxpool.Pool
	salt string
}

func main() {
	flag.Parse()
	ctx := context.Background()

	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	hashingConfig, err := config.NewHashingConfig()
	if err != nil {
		log.Fatalf("failed to load hashing config: %v", err)
	}

	grpcConfig, err := config.NewGRPCConfig()
	if err != nil {
		log.Fatalf("failed to get grpc config: %v", err)
	}

	pgConfig, err := config.NewPGConfig()
	if err != nil {
		log.Fatalf("failed to get pg config: %v", err)
	}

	lis, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	pool, err := pgxpool.Connect(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterAuthV1Server(s, &server{pool: pool, salt: hashingConfig.Salt()})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// Create user
func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	log.Printf("User info: %v, password: %v, password confirm: %v", req.GetInfo(), req.GetPassword(), req.GetPasswordConfirm())

	if req.Info == nil {
		return nil, status.Error(codes.InvalidArgument, "info is required")
	}

	if req.Info.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	if req.Info.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	if req.Info.Role == 0 {
		return nil, status.Error(codes.InvalidArgument, "role is required")
	}

	if req.Password != req.PasswordConfirm {
		return nil, status.Error(codes.InvalidArgument, "passwords do not match")
	}

	password, err := hashing.HashPasswordWithSalt(req.Password, s.salt)
	if err != nil {
		return nil, status.Error(codes.Internal, "could not hash password")
	}

	builderInsert := sq.Insert("\"user\"").
		PlaceholderFormat(sq.Dollar).
		Columns("name", "email", "role", "password").
		Values(req.Info.Name, req.Info.Email, req.Info.Role, password).
		Suffix("RETURNING user_id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	var userID int
	err = s.pool.QueryRow(ctx, query, args...).Scan(&userID)
	if err != nil {
		log.Fatalf("failed to insert user: %v", err)
	}

	log.Printf("inserted user with id: %d", userID)

	return &desc.CreateResponse{
		Id: int64(userID),
	}, nil
}

// Get user by id
func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
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

	var id int64
	var name, email string
	var role int32
	var createdAt time.Time
	var updatedAt sql.NullTime

	err = s.pool.QueryRow(ctx, query, args...).Scan(&id, &name, &email, &role, &createdAt, &updatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		log.Fatalf("failed to select users: %v", err)
	}

	var outUpdatedAt *timestamppb.Timestamp
	if updatedAt.Valid {
		outUpdatedAt = timestamppb.New(updatedAt.Time)
	}

	return &desc.GetResponse{
		User: &desc.User{
			Id: id,
			Info: &desc.UserInfo{
				Name:  name,
				Email: email,
				Role:  desc.Role(role),
			},
			CreatedAt: timestamppb.New(createdAt),
			UpdatedAt: outUpdatedAt,
		},
	}, nil
}

// Update user by id
func (s *server) Update(ctx context.Context, req *desc.UpdateRequest) (*empty.Empty, error) {
	log.Printf("User info: %+v", req.GetInfo())

	info := req.GetInfo()
	if info == nil {
		return nil, status.Error(codes.InvalidArgument, "info is required")
	}

	builderUpdate := sq.Update("\"user\"").
		PlaceholderFormat(sq.Dollar)

	if info.GetName() != nil {
		builderUpdate = builderUpdate.Set("name", info.GetName().GetValue())
	}

	if info.GetRole() != 0 {
		builderUpdate = builderUpdate.Set("role", info.GetRole().Number())
	}

	builderUpdate = builderUpdate.Set("updated_at", time.Now())
	builderUpdate = builderUpdate.Where(sq.Eq{"user_id": req.GetId()})

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	res, err := s.pool.Exec(ctx, query, args...)
	if err != nil {
		log.Fatalf("failed to update user: %v", err)
	}

	log.Printf("updated %d rows", res.RowsAffected())

	return nil, nil
}

// Delete user by id
func (s *server) Delete(ctx context.Context, req *desc.DeleteRequest) (*empty.Empty, error) {
	log.Printf("User id: %d", req.GetId())

	builderDelete := sq.Delete("\"user\"").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"user_id": req.GetId()})

	query, args, err := builderDelete.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	res, err := s.pool.Exec(ctx, query, args...)
	if err != nil {
		log.Fatalf("failed to delete user: %v", err)
	}

	log.Printf("delete %d rows", res.RowsAffected())

	return nil, nil
}
