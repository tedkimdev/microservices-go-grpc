package main

import (
  "fmt"
  "github.com/tedkimdev/microservices-go-grpc/pkg/database"
  "log"
  "net"

  "github.com/tedkimdev/microservices-go-grpc/internal/auth"
  config "github.com/tedkimdev/microservices-go-grpc/pkg/config/auth_svc"
  "github.com/tedkimdev/microservices-go-grpc/pkg/utils"
  pb "github.com/tedkimdev/microservices-go-grpc/proto/v1/auth"
  "google.golang.org/grpc"
)

func main() {
  c, err := config.LoadConfig()
  if err != nil {
    log.Fatalln("Failed at config", err)
  }

  db, err := database.NewDatabase(c.DBUrl)
  if err != nil {
    log.Fatalln("Failed to connect database", err)
  }

  jwt := utils.JwtWrapper{
    SecretKey:       c.JWTSecretKey,
    Issuer:          "go-grpc-auth-svc",
    ExpirationHours: 24 * 365,
  }

  lis, err := net.Listen("tcp", c.Port)

  if err != nil {
    log.Fatalln("Failed to listing:", err)
  }

  fmt.Println("Auth Svc on", c.Port)

  s := auth.NewServiceServer(db, jwt)

  grpcServer := grpc.NewServer()

  pb.RegisterAuthServiceServer(grpcServer, s)

  if err := grpcServer.Serve(lis); err != nil {
    log.Fatalln("Failed to serve:", err)
  }
}