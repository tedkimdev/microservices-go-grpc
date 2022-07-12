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

func (s *ServiceServer) CreateProduct(_ context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
  var product models.Product

  product.Name = req.Name
  product.Stock = req.Stock
  product.Price = req.Price

  if result := s.cfg.ReadWriteDB().Create(&product); result.Error != nil {
    return &pb.CreateProductResponse{
      Status: http.StatusConflict,
      Error:  result.Error.Error(),
    }, nil
  }

  return &pb.CreateProductResponse{
    Status: http.StatusCreated,
    Id:     product.Id,
  }, nil
}

func (s *ServiceServer) FindOne(_ context.Context, req *pb.FindOneRequest) (*pb.FindOneResponse, error) {
  var product models.Product

  if result := s.cfg.ReadOnlyDB().First(&product, req.Id); result.Error != nil {
    return &pb.FindOneResponse{
      Status: http.StatusNotFound,
      Error:  result.Error.Error(),
    }, nil
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

func (s *ServiceServer) DecreaseStock(_ context.Context, req *pb.DecreaseStockRequest) (*pb.DecreaseStockResponse, error) {
  var product models.Product

  if result := s.cfg.ReadOnlyDB().First(&product, req.Id); result.Error != nil {
    return &pb.DecreaseStockResponse{
      Status: http.StatusNotFound,
      Error:  result.Error.Error(),
    }, nil
  }

  if product.Stock <= 0 {
    return &pb.DecreaseStockResponse{
      Status: http.StatusConflict,
      Error:  "Stock too low",
    }, nil
  }

  var log models.StockDecreaseLog

  if result := s.cfg.ReadOnlyDB().Where(&models.StockDecreaseLog{OrderId: req.OrderId}).First(&log); result.Error == nil {
    return &pb.DecreaseStockResponse{
      Status: http.StatusConflict,
      Error:  "Stock already decreased",
    }, nil
  }

  product.Stock = product.Stock - 1

  s.cfg.ReadWriteDB().Save(&product)

  log.OrderId = req.OrderId
  log.ProductRefer = product.Id

  s.cfg.ReadWriteDB().Create(&log)

  return &pb.DecreaseStockResponse{
    Status: http.StatusOK,
  }, nil
}