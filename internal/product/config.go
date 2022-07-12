package product

import "gorm.io/gorm"

type Config interface {
  ReadOnlyDB() *gorm.DB
  ReadWriteDB() *gorm.DB
  ProductRepo() Repository
}

type config struct {
  Config
  readOnlyDB *gorm.DB
  readWriteDB *gorm.DB
  productRepo Repository
}

func NewConfig(readOnlyDB *gorm.DB, readWriteDB *gorm.DB, productRepo Repository) Config {
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
func (c *config) ProductRepo() Repository {
  return c.productRepo
}