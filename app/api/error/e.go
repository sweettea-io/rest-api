package error

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
    Data: utils.JSON{
      "ok": true,
      "code": code,
      "error": msg,
    },
  }
}

func Forbidden() (*Error) {
  return ApiError(http.StatusForbidden, http.StatusForbidden, "forbidden")
}

func Unauthorized() (*Error) {
  return ApiError(http.StatusUnauthorized, http.StatusUnauthorized, "unauthorized")
}

func UnknownError() (*Error) {
  return ApiError(http.StatusInternalServerError, http.StatusInternalServerError, "unknown_error")
}

func InvalidPayload() (*Error) {
  return ApiError(http.StatusBadRequest, http.StatusBadRequest, "invalid_input_payload")
}