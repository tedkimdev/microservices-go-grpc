package handler

import (
  "context"
  "github.com/tedkimdev/microservices-go-grpc/internal/repo"
  "github.com/tedkimdev/microservices-go-grpc/pkg/utils"
  pb "github.com/tedkimdev/microservices-go-grpc/proto/v1/auth"
  "net/http"
)

type LoginFunc func(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error)

func Login(authRepo repo.AuthRepository, jwt *utils.JwtWrapper) LoginFunc {
  return func(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
    user, err := authRepo.GetUserByEmail(ctx, req.Email)
    if err != nil {
      return &pb.LoginResponse{
        Status: http.StatusNotFound,
        Error:  "User not found",
      }, nil
    }

    match := utils.CheckPasswordHash(req.Password, user.Password)

    if !match {
      return &pb.LoginResponse{
        Status: http.StatusNotFound,
        Error:  "User not found",
      }, nil
    }

    token, _ := jwt.GenerateToken(*user)

    return &pb.LoginResponse{
      Status: http.StatusOK,
      Token:  token,
    }, nil
  }
}