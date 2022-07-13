package auth

import (
  "context"
  "github.com/tedkimdev/microservices-go-grpc/internal/auth/handler"
  "github.com/tedkimdev/microservices-go-grpc/pkg/utils"
  pb "github.com/tedkimdev/microservices-go-grpc/proto/v1/auth"
  "google.golang.org/grpc"
  "gorm.io/gorm"
)

type ServiceServer struct {
  pb.AuthServiceServer

  cfg Config
  Jwt *utils.JwtWrapper
  readOnlyDB *gorm.DB
}

func NewServiceServer(cfg Config) (*ServiceServer, error) {
  return &ServiceServer{
    cfg: cfg,
    Jwt: cfg.Jwt(),
    readOnlyDB: cfg.ReadOnlyDB(),
  }, nil
}

func NewGRPCServer(cfg Config) (*grpc.Server, error) {
  authServer, err := NewServiceServer(cfg)
  if err != nil {
    return nil, err
  }
  grpcServer := grpc.NewServer()
  pb.RegisterAuthServiceServer(grpcServer, authServer)
  return grpcServer, nil
}


func (s *ServiceServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
  return handler.Register(s.cfg.AuthRepo())(ctx, req)
}

func (s *ServiceServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
  return handler.Login(s.cfg.AuthRepo(), s.cfg.Jwt())(ctx, req)
}

func (s *ServiceServer) Validate(ctx context.Context, req *pb.ValidateRequest) (*pb.ValidateResponse, error) {
  return handler.Validate(s.cfg.AuthRepo(), s.cfg.Jwt())(ctx, req)
}