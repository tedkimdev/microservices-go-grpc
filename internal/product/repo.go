package product

import (
  "gorm.io/gorm"
)

type Repository interface {
}

type repository struct {
  readOnlyDB *gorm.DB
  readWriteDB *gorm.DB
}

func NewRepository(
  readOnlyDB *gorm.DB,
  readWriteDB *gorm.DB,
) Repository {
  return &repository{
    readOnlyDB:  readOnlyDB,
    readWriteDB: readWriteDB,
  }
}