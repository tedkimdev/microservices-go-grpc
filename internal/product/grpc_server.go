package product

import (
  "context"
  "github.com/tedkimdev/microservices-go-grpc/pkg/models"
  pb "github.com/tedkimdev/microservices-go-grpc/proto/v1/product"
  "google.golang.org/grpc"
  "gorm.io/gorm"
  "net/http"
)

type ServiceServer struct {
  pb.ProductServiceServer

  cfg Config
  readOnlyDB *gorm.DB
}

func NewServiceServer(cfg Config) (*ServiceServer, error) {
  return &ServiceServer{
    cfg: cfg,
    readOnlyDB: cfg.ReadOnlyDB(),
  }, nil
}
func NewGRPCServer(cfg Config) (*grpc.Server, error) {
  productServer, err := NewServiceServer(cfg)
  if err != nil {
    return nil, err
  }
  grpcServer := grpc.NewServer()
  pb.RegisterProductServiceServer(grpcServer, productServer)
  return grpcServer, nil
  //logrus.ErrorKey = "grpc.error"
  //logrusEntry := logrus.NewEntry(logrus.StandardLogger())
  //statsdClient := cfg.StatsdClient().CloneWithPrefixExtension(".grpc")
  //
  //grpcServer := grpc.NewServer(
  //  grpc_middleware.WithUnaryServerChain(
  //    grpc_ctxtags.UnaryServerInterceptor(
  //      grpc_ctxtags.WithFieldExtractor(
  //        grpc_ctxtags.CodeGenRequestFieldExtractor,
  //      ),
  //    ),
  //    grpc_logrus.UnaryServerInterceptor(logrusEntry),
  //    banksalad.StatsUnaryServerInterceptor(statsdClient),
  //    banksalad.LogStackTraceUnaryServerInterceptor(),
  //    banksalad.CanceledErrorUnaryServerInterceptor(),
  //    banksalad.LogCallerUnaryServerInterceptor(),
  //    grpc_recovery.UnaryServerInterceptor(
  //      grpc_recovery.WithRecoveryHandlerContext(banksalad.LogPanicStackTrace()),
  //    ),
  //  ),
  //  grpc.KeepaliveParams(keepalive.ServerParameters{
  //    MaxConnectionAge: 30 * time.Second,
  //  }),
  //)
  //
  //loancurationServer, err := NewLoancurationServer(cfg)
  //if err != nil {
  //  return nil, err
  //}
  //
  //loancuration.RegisterLoancurationServer(grpcServer, loancurationServer)
  //reflection.Register(grpcServer)
  //
  //return grpcServer, nil
}

func (s *ServiceServer) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
  var product models.Product

  product.Name = req.Name
  product.Stock = req.Stock
  product.Price = req.Price

  newProduct, err := s.cfg.ProductRepo().CreateProduct(ctx, &product)
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

func (s *ServiceServer) FindOne(ctx context.Context, req *pb.FindOneRequest) (*pb.FindOneResponse, error) {

  productID := req.Id
  product, err := s.cfg.ProductRepo().GetProductByID(ctx, productID)
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

func (s *ServiceServer) DecreaseStock(ctx context.Context, req *pb.DecreaseStockRequest) (*pb.DecreaseStockResponse, error) {
  product, err := s.cfg.ProductRepo().GetProductByID(ctx, req.Id)
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

  log, err := s.cfg.ProductRepo().GetStockLog(ctx, req.OrderId)
  if err != nil {
    return &pb.DecreaseStockResponse{
      Status: http.StatusInternalServerError,
      Error:  "Stock already decreased",
    }, err
  }

  product.Stock = product.Stock - 1

  _, err = s.cfg.ProductRepo().UpdateProduct(ctx, product)
  if err != nil {
    return &pb.DecreaseStockResponse{
      Status: http.StatusInternalServerError,
      Error:  "Failed to update product",
    }, err
  }

  log.OrderId = req.OrderId
  log.ProductRefer = product.Id
  _, err = s.cfg.ProductRepo().UpdateStockLog(ctx, log)
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