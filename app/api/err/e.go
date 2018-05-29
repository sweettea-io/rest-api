package err

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

func Forbidden() (*Error) {
  status := http.StatusForbidden
  return ApiError(status, status, http.StatusText(status))
}

func Unauthorized() (*Error) {
  status := http.StatusUnauthorized
  return ApiError(status, status, http.StatusText(status))
}

func UnknownError() (*Error) {
  return ApiError(http.StatusInternalServerError, http.StatusInternalServerError, "unknown_error")
}

func InvalidPayload() (*Error) {
  return ApiError(http.StatusBadRequest, http.StatusBadRequest, "invalid_input_payload")
}