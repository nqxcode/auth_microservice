package main

import (
	"context"
	"fmt"
	"log"
	"net"

	desc "github.com/nqxcode/auth_microservice/pkg/auth_v1"

	"github.com/brianvoe/gofakeit"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const grpcPort = 50051

type server struct {
	desc.UnimplementedAuthV1Server
}

// Create user
func (s *server) Create(_ context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	log.Printf("User info: %v, password: %v, password confirm: %v", req.GetInfo(), req.GetPassword(), req.GetPasswordConfirm())

	return &desc.CreateResponse{
		Id: int64(gofakeit.Number(1, 1000)),
	}, nil
}

// Get user by id
func (s *server) Get(_ context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	log.Printf("User id: %d", req.GetId())

	return &desc.GetResponse{
		User: &desc.User{
			Id: req.GetId(),
			Info: &desc.UserInfo{
				Name:  gofakeit.Name(),
				Email: gofakeit.Email(),
				Role: func() desc.Role {
					if gofakeit.Number(0, 1) == 0 {
						return desc.Role_ADMIN
					}

					return desc.Role_USER
				}(),
			},
			CreatedAt: timestamppb.New(gofakeit.Date()),
			UpdatedAt: timestamppb.New(gofakeit.Date()),
		},
	}, nil
}

// Update user by id
func (s *server) Update(_ context.Context, req *desc.UpdateRequest) (*empty.Empty, error) {
	log.Printf("User info: %+v", req.GetInfo())

	return nil, nil
}

// Delete user by id
func (s *server) Delete(_ context.Context, req *desc.DeleteRequest) (*empty.Empty, error) {
	log.Printf("User id: %d", req.GetId())

	return nil, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterAuthV1Server(s, &server{})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
