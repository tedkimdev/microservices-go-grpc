package handler

import (
  "context"
  "github.com/tedkimdev/microservices-go-grpc/internal/repo"
  pb "github.com/tedkimdev/microservices-go-grpc/proto/v1/product"
  "net/http"
)

type DecreaseStockFunc func(ctx context.Context, req *pb.DecreaseStockRequest) (*pb.DecreaseStockResponse, error)

func DecreaseStock(productRepo repo.ProductRepository) DecreaseStockFunc {
  return func(ctx context.Context, req *pb.DecreaseStockRequest) (*pb.DecreaseStockResponse, error) {
    product, err := productRepo.GetProductByID(ctx, req.Id)
    if err != nil {
      return &pb.DecreaseStockResponse{
        Status: http.StatusNotFound,
        Error:  err.Error(),
      }, err
    }

    if product.Stock <= 0 {
      return &pb.DecreaseStockResponse{
        Status: http.StatusConflict,
        Error:  "Stock too low",
      }, nil
    }

    log, err := productRepo.GetStockLog(ctx, req.OrderId)
    if err != nil {
      return &pb.DecreaseStockResponse{
        Status: http.StatusInternalServerError,
        Error:  "Stock already decreased",
      }, err
    }

    product.Stock = product.Stock - 1

    _, err = productRepo.UpdateProduct(ctx, product)
    if err != nil {
      return &pb.DecreaseStockResponse{
        Status: http.StatusInternalServerError,
        Error:  "Failed to update product",
      }, err
    }

    log.OrderId = req.OrderId
    log.ProductRefer = product.Id
    _, err = productRepo.UpdateStockLog(ctx, log)
    if err != nil {
      return &pb.DecreaseStockResponse{
        Status: http.StatusInternalServerError,
        Error:  "Failed to update StockDecreaseLog",
      }, err
    }

    return &pb.DecreaseStockResponse{
      Status: http.StatusOK,
    }, nil
  }
}