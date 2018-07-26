package sessionsvc

import (
  "fmt"
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/pkg/model"
)

func Create(user *model.User) (*model.Session, error) {
  // Create Session model.
  session := model.Session{User: user}

  // Create record
  if err := app.DB.Create(&session).Error; err != nil {
    return nil, fmt.Errorf("error creating Session: %s", err.Error())
  }

  return &session, nil
}
