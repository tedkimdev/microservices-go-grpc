package repo

import (
  "context"
  "github.com/tedkimdev/microservices-go-grpc/pkg/models"
  "gorm.io/gorm"
)

type AuthRepository interface {
  GetUserByEmail(_ context.Context, email string) (*models.User, error)
  CreateUser(_ context.Context, user *models.User) (*models.User, error)
}

type authRepository struct {
  readOnlyDB *gorm.DB
  readWriteDB *gorm.DB
}

func NewAuthRepository(
  readOnlyDB *gorm.DB,
  readWriteDB *gorm.DB,
) AuthRepository {
  return &authRepository{
    readOnlyDB:  readOnlyDB,
    readWriteDB: readWriteDB,
  }
}

func (r *authRepository) GetUserByEmail(_ context.Context, email string) (*models.User, error) {
  var user models.User
  if result := r.readOnlyDB.Where(&models.User{Email:email}).First(&user); result.Error != nil {
    return nil, result.Error
  }
  return &user, nil
}

func (r *authRepository) CreateUser(_ context.Context, user *models.User) (*models.User, error) {
  if result := r.readWriteDB.Create(user); result.Error != nil {
    return nil, result.Error
  }
  return user, nil
}