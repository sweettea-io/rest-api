package usersvc

import (
  "fmt"
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/pkg/model"
)

func Create(email string, hashedPw string, admin bool) (*model.User, error) {
  // Create User model.
  user := model.User{
    Email: email,
    HashedPw: hashedPw,
    Admin: admin,
  }

  // Create record.
  if err := app.DB.Create(&user).Error; err != nil {
    return nil, fmt.Errorf("error creating User: %s", err.Error())
  }

  return &user, nil
}