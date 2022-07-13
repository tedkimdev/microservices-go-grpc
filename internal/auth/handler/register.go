package handler

import (
  "context"
  "github.com/tedkimdev/microservices-go-grpc/internal/repo"
  "github.com/tedkimdev/microservices-go-grpc/pkg/utils"
  pb "github.com/tedkimdev/microservices-go-grpc/proto/v1/auth"
  "net/http"
)

type RegisterFunc func(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error)

func Register(authRepo repo.AuthRepository) RegisterFunc {
  return func(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
    user, err := authRepo.GetUserByEmail(ctx, req.Email)
    if err != nil {
      return &pb.RegisterResponse{
        Status: http.StatusConflict,
        Error:  "E-Mail already exists",
      }, nil
    }

    user.Email = req.Email
    user.Password = utils.HashPassword(req.Password)

    _, err = authRepo.CreateUser(ctx, user)
    if err != nil {
      return &pb.RegisterResponse{
        Status: http.StatusInternalServerError,
        Error:  "Failed to create user",
      }, nil
    }

    return &pb.RegisterResponse{
      Status: http.StatusCreated,
    }, nil
  }
}
