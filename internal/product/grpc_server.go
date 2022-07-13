package product

import (
  "context"
  "github.com/tedkimdev/microservices-go-grpc/internal/product/handler"
  pb "github.com/tedkimdev/microservices-go-grpc/proto/v1/product"
  "google.golang.org/grpc"
  "gorm.io/gorm"
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
  return handler.CreateProduct(s.cfg.ProductRepo())(ctx, req)
}

func (s *ServiceServer) FindOne(ctx context.Context, req *pb.FindOneRequest) (*pb.FindOneResponse, error) {
  return handler.FindOneProduct(s.cfg.ProductRepo())(ctx, req)
}

func (s *ServiceServer) DecreaseStock(ctx context.Context, req *pb.DecreaseStockRequest) (*pb.DecreaseStockResponse, error) {
  return handler.DecreaseStock(s.cfg.ProductRepo())(ctx, req)
}