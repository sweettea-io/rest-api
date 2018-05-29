package api

import (
  "net/http"
  "github.com/sweettea/rest-api/defs"
  "github.com/sweettea/rest-api/pkg/models"
  "errors"
)

func LoadCurrentUser(w http.ResponseWriter, req *http.Request, user *models.User) error {
  token := req.Header.Get(defs.AuthHeaderName)

  if token == "" {
    return errors.New(http.StatusText(http.StatusUnauthorized))
  }

  // Find session by token
  var session models.Session
  db.Where(&models.Session{Token: token}).First(&session)

  // Load reference to user
  user = &session.User

  return nil
}