package api

import (
  "net/http"
  "github.com/sweettea-io/rest-api/defs"
  "github.com/sweettea-io/rest-api/pkg/models"
  "errors"
)

// Find and set current user from session token provided in request header.
func LoadCurrentUser(w http.ResponseWriter, req *http.Request, user *models.User) error {
  token := req.Header.Get(defs.AuthHeaderName)

  // Return error if auth header token not found.
  if token == "" {
    return errors.New(http.StatusText(http.StatusUnauthorized))
  }

  // Find session by token.
  var session models.Session
  result := db.Preload("User", "is_destroyed = ?", false).
    Where(&models.Session{Token: token}).
    First(&session)

  // Return error if session doesn't exist or user doesn't exist through session.
  if result.RecordNotFound() || session.User.ID == 0 {
    return errors.New(http.StatusText(http.StatusUnauthorized))
  }

  // Update user address.
  user = &session.User

  return nil
}