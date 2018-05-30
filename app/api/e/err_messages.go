package e

import (
  "net/http"
  "github.com/sweettea/rest-api/pkg/utils"
)

type Error struct {
  Status int
  Data utils.JSON
}

func ApiError(status int, code int, msg string) (*Error) {
  return &Error{
    Status: status,
    Data: utils.JSON{"ok": true, "code": code, "error": msg},
  }
}

// ------- Generic Errors --------

const JsonEncodingError = `{"ok": false, "code": 1000, "error": "Error encoding JSON response"}`

func Forbidden() (*Error) {
  status := http.StatusForbidden
  return ApiError(status, status, http.StatusText(status))
}

func Unauthorized() (*Error) {
  status := http.StatusUnauthorized
  return ApiError(status, status, http.StatusText(status))
}

func ISE() (*Error) {
  status := http.StatusInternalServerError
  return ApiError(status, status, http.StatusText(status))
}

func InvalidPayload() (*Error) {
  return ApiError(http.StatusBadRequest, http.StatusBadRequest, "invalid_input_payload")
}

// ------- User Errors --------

func UserNotFound() (*Error) {
  return ApiError(http.StatusNotFound, 1001, "user_not_found")
}