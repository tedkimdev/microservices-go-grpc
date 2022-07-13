package auth

import (
  "github.com/tedkimdev/microservices-go-grpc/internal/repo"
  "github.com/tedkimdev/microservices-go-grpc/pkg/utils"
  "gorm.io/gorm"
)

type Config interface {
  ReadOnlyDB() *gorm.DB
  ReadWriteDB() *gorm.DB
  AuthRepo() repo.AuthRepository
  Jwt() *utils.JwtWrapper
}

type config struct {
  Config
  readOnlyDB *gorm.DB
  readWriteDB *gorm.DB
  authRepo repo.AuthRepository
  jwt *utils.JwtWrapper
}

func NewConfig(readOnlyDB *gorm.DB, readWriteDB *gorm.DB, authRepo repo.AuthRepository, jwt *utils.JwtWrapper) Config {
  return &config{
    readOnlyDB: readOnlyDB,
    readWriteDB: readWriteDB,
    authRepo: authRepo,
    jwt: jwt,
  }
}

func (c *config) ReadOnlyDB() *gorm.DB {
  return c.readOnlyDB
}

func (c *config) ReadWriteDB() *gorm.DB {
  return c.readWriteDB
}

func (c *config) AuthRepo() repo.AuthRepository {
  return c.authRepo
}

func (c *config) Jwt() *utils.JwtWrapper {
  return c.jwt
}