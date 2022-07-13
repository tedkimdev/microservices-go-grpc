package handler

import (
  "context"
  "github.com/tedkimdev/microservices-go-grpc/internal/repo"
  pb "github.com/tedkimdev/microservices-go-grpc/proto/v1/product"
  "net/http"
)

type FindOneFunc func(ctx context.Context, req *pb.FindOneRequest) (*pb.FindOneResponse, error)

func FindOneProduct(productRepo repo.ProductRepository) FindOneFunc {
  return func(ctx context.Context, req *pb.FindOneRequest) (*pb.FindOneResponse, error) {
    productID := req.Id
    product, err := productRepo.GetProductByID(ctx, productID)
    if err != nil {
      return &pb.FindOneResponse{
        Status: http.StatusNotFound,
        Error:  err.Error(),
      }, err
    }

    data := &pb.FindOneData{
      Id:    product.Id,
      Name:  product.Name,
      Stock: product.Stock,
      Price: product.Price,
    }

    return &pb.FindOneResponse{
      Status: http.StatusOK,
      Data:   data,
    }, nil
  }
}