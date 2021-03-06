package usersvc

import (
  "errors"
  "fmt"
  "net/http"
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/pkg/model"
)

func FromRequest(req *http.Request) (*model.User, error) {
  // Get auth token from request header.
  token := req.Header.Get(app.Config.AuthHeaderName)

  // Error to be returned if token or user isn't found.
  e := errors.New(http.StatusText(http.StatusUnauthorized))

  // Return error if auth header token not found.
  if token == "" {
    return nil, e
  }

  // Find session by token.
  var session model.Session
  result := app.DB.Preload("User").Where(&model.Session{Token: token}).First(&session)

  // Return error if session doesn't exist or user doesn't exist through session.
  if result.RecordNotFound() || session.User.ID == 0 {
    return nil, e
  }

  return session.User, nil
}

func FromEmail(email string) (*model.User, error) {
  // Find User by email.
  var user model.User
  result := app.DB.Where(&model.User{Email: email}).Find(&user)

  // Return error if not found.
  if result.RecordNotFound() {
    return nil, fmt.Errorf("User(email=%s) not found.\n", email)
  }

  return &user, nil
}