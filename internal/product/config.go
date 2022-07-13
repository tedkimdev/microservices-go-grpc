package product

import (
  "github.com/tedkimdev/microservices-go-grpc/internal/repo"
  "gorm.io/gorm"
)

type Config interface {
  ReadOnlyDB() *gorm.DB
  ReadWriteDB() *gorm.DB
  ProductRepo() repo.ProductRepository
}

type config struct {
  Config
  readOnlyDB  *gorm.DB
  readWriteDB *gorm.DB
  productRepo repo.ProductRepository
}

func NewConfig(readOnlyDB *gorm.DB, readWriteDB *gorm.DB, productRepo repo.ProductRepository) Config {
  return &config{
    readOnlyDB: readOnlyDB,
    readWriteDB: readWriteDB,
    productRepo: productRepo,
  }
}

func (c *config) ReadOnlyDB() *gorm.DB {
  return c.readOnlyDB
}
func (c *config) ReadWriteDB() *gorm.DB {
  return c.readWriteDB
}
func (c *config) ProductRepo() repo.ProductRepository {
  return c.productRepo
}