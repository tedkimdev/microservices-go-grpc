package handler

import (
  "context"
  "github.com/tedkimdev/microservices-go-grpc/internal/repo"
  "github.com/tedkimdev/microservices-go-grpc/pkg/models"
  pb "github.com/tedkimdev/microservices-go-grpc/proto/v1/product"
  "net/http"
)

type CreateProductFunc func(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error)

func CreateProduct(productRepo repo.ProductRepository) CreateProductFunc {
  return func(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
    var product models.Product

    product.Name = req.Name
    product.Stock = req.Stock
    product.Price = req.Price

    newProduct, err := productRepo.CreateProduct(ctx, &product)
    if err != nil {
      return &pb.CreateProductResponse{
        Status: http.StatusConflict,
        Error:  err.Error(),
      }, err
    }

    return &pb.CreateProductResponse{
      Status: http.StatusCreated,
      Id:     newProduct.Id,
    }, nil
  }
}