package handler

import (
  "context"
  "github.com/tedkimdev/microservices-go-grpc/internal/repo"
  "github.com/tedkimdev/microservices-go-grpc/pkg/utils"
  pb "github.com/tedkimdev/microservices-go-grpc/proto/v1/auth"
  "net/http"
)

type ValidateFunc func(ctx context.Context, req *pb.ValidateRequest) (*pb.ValidateResponse, error)

func Validate(authRepo repo.AuthRepository, jwt *utils.JwtWrapper) ValidateFunc {
  return func(ctx context.Context, req *pb.ValidateRequest) (*pb.ValidateResponse, error) {
    claims, err := jwt.ValidateToken(req.Token)

    if err != nil {
      return &pb.ValidateResponse{
        Status: http.StatusBadRequest,
        Error:  err.Error(),
      }, nil
    }

    user, err := authRepo.GetUserByEmail(ctx, claims.Email)
    if err != nil {
      return &pb.ValidateResponse{
        Status: http.StatusNotFound,
        Error:  "User not found",
      }, nil
    }

    return &pb.ValidateResponse{
      Status: http.StatusOK,
      UserId: user.Id,
    }, nil
  }
}