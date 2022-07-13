package repo

import (
  "context"
  "github.com/tedkimdev/microservices-go-grpc/pkg/models"
  "gorm.io/gorm"
)

type ProductRepository interface {
  CreateProduct(_ context.Context, product *models.Product) (*models.Product, error)
  GetProductByID(_ context.Context, id int64) (*models.Product, error)
  UpdateProduct(_ context.Context, product *models.Product) (*models.Product, error)

  GetStockLog(_ context.Context, orderId int64) (*models.StockDecreaseLog, error)
  UpdateStockLog(_ context.Context, stockLog *models.StockDecreaseLog) (*models.StockDecreaseLog, error)
}

type productRepository struct {
  readOnlyDB *gorm.DB
  readWriteDB *gorm.DB
}

func NewRepository(
  readOnlyDB *gorm.DB,
  readWriteDB *gorm.DB,
) ProductRepository {
  return &productRepository{
    readOnlyDB:  readOnlyDB,
    readWriteDB: readWriteDB,
  }
}

func (r *productRepository) GetStockLog(_ context.Context, orderId int64) (*models.StockDecreaseLog, error) {
  var log models.StockDecreaseLog
  if result := r.readOnlyDB.Where(&models.StockDecreaseLog{OrderId: orderId}).First(&log); result.Error == nil {
    return nil, result.Error
  }
  return &log, nil
}

func (r *productRepository) CreateProduct(_ context.Context, product *models.Product) (*models.Product, error) {
  if result := r.readWriteDB.Create(&product); result.Error != nil {
    return nil, result.Error
  }
  return product, nil
}

func (r *productRepository) GetProductByID(_ context.Context, id int64) (*models.Product, error) {
  var product models.Product
  if result := r.readOnlyDB.First(&product, id); result.Error != nil {
    return nil, result.Error
  }
  return &product, nil
}

func (r *productRepository) UpdateProduct(_ context.Context, product *models.Product) (*models.Product, error) {
  if result := r.readWriteDB.Save(product); result.Error != nil {
    return nil, result.Error
  }
  return product, nil
}

func (r *productRepository) UpdateStockLog(_ context.Context, stockLog *models.StockDecreaseLog) (*models.StockDecreaseLog, error) {
  if result := r.readWriteDB.Save(stockLog); result.Error != nil {
    return nil, result.Error
  }
  return stockLog, nil
}