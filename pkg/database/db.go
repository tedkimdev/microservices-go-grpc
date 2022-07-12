package database

import (
  "fmt"
  "github.com/tedkimdev/microservices-go-grpc/pkg/models"
  "gorm.io/driver/postgres"
  "gorm.io/gorm"
)

func NewDatabase(url string) (*gorm.DB, error) {
  db, err := gorm.Open(postgres.Open(url), &gorm.Config{})

  if err != nil {
    return nil, err
  }

  db.AutoMigrate(&models.User{})

  return db, nil
}

func NewProductDatabase(url string) (*gorm.DB, error) {
  fmt.Println(url)
  db, err := gorm.Open(postgres.Open(url), &gorm.Config{})

  if err != nil {
    return nil, err
  }

  db.AutoMigrate(&models.Product{})
  db.AutoMigrate(&models.StockDecreaseLog{})

  return db, nil
}