package main

import (
  "fmt"
  "github.com/tedkimdev/microservices-go-grpc/internal/repo"
  "github.com/tedkimdev/microservices-go-grpc/pkg/database"
  "log"
  "net"

  "github.com/tedkimdev/microservices-go-grpc/internal/auth"
  config "github.com/tedkimdev/microservices-go-grpc/pkg/config/auth_svc"
  "github.com/tedkimdev/microservices-go-grpc/pkg/utils"
)

func main() {
  envConfig, err := config.LoadEnvConfig()
  if err != nil {
    log.Fatalln("Failed at config", err)
  }

  db, err := database.NewDatabase(envConfig.ReadOnlyDBURL)
  if err != nil {
    log.Fatalln("Failed to connect database", err)
  }

  jwt := utils.JwtWrapper{
    SecretKey:       envConfig.JWTSecretKey,
    Issuer:          "go-grpc-auth-svc",
    ExpirationHours: 24 * 365,
  }

  lis, err := net.Listen("tcp", envConfig.Port)

  if err != nil {
    log.Fatalln("Failed to listing:", err)
  }
  fmt.Println("Auth Svc on", envConfig.Port)

  ar := repo.NewAuthRepository(db, nil)
  cfg := auth.NewConfig(db, nil, ar, &jwt)
  grpcServer, err := auth.NewGRPCServer(cfg)

  if err := grpcServer.Serve(lis); err != nil {
    log.Fatalln("Failed to serve:", err)
  }
}