package database

import (
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