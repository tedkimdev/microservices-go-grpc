package main

import (
  "fmt"
  "github.com/tedkimdev/microservices-go-grpc/internal/repo"
  "log"
  "net"

  "github.com/tedkimdev/microservices-go-grpc/internal/product"
  config "github.com/tedkimdev/microservices-go-grpc/pkg/config/product_svc"
  "github.com/tedkimdev/microservices-go-grpc/pkg/database"
)

func main() {
  envConfig, err := config.LoadEnvConfig()
  if err != nil {
    log.Fatalln("Failed at config", err)
  }

  readOnlyDB, err := database.NewProductDatabase(envConfig.ReadOnlyDBURL)
  if err != nil {
    log.Fatalln("Failed to connect readonly database", err)
  }
  //readWriteDB, err := database.NewProductDatabase(envConfig.ReadWriteDBURL)
  //if err != nil {
  // log.Fatalln("Failed to connect readwrite database", err)
  //}

  lis, err := net.Listen("tcp", envConfig.Port)
  if err != nil {
    log.Fatalln("Failed to listing:", err)
  }
  fmt.Println("Product Service on", envConfig.Port)

  pr := repo.NewProductRepository(readOnlyDB, nil)
  cfg := product.NewConfig(readOnlyDB, nil, pr)
  grpcServer, err := product.NewGRPCServer(cfg)

  if err := grpcServer.Serve(lis); err != nil {
    log.Fatalln("Failed to serve:", err)
  }
}