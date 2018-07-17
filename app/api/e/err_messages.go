package e

import (
  "net/http"
  "github.com/sweettea-io/rest-api/pkg/utils"
)

type Error struct {
  Status int
  Data utils.JSON
}

func ApiError(status int, code int, msg string) (*Error) {
  return &Error{
    Status: status,
    Data: utils.JSON{
      "ok": false,
      "code": code,
      "error": msg,
    },
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

func EmailNotAvailable() (*Error) {
  return ApiError(http.StatusInternalServerError, 1002, "email_not_available")
}

func UserCreationFailed() (*Error) {
  return ApiError(http.StatusInternalServerError, 1003, "user_creation_failed")
}

// ------- Session Errors --------

func SessionCreationFailed() (*Error) {
  return ApiError(http.StatusInternalServerError, 2001, "session_creation_failed")
}

// ------- Company Errors --------

func CompanyAlreadyExists() (*Error) {
  return ApiError(http.StatusInternalServerError, 3001, "company_already_exists")
}

func CompanyCreationFailed() (*Error) {
  return ApiError(http.StatusInternalServerError, 3002, "company_creation_failed")
}

func CompanyNotFound() (*Error) {
  return ApiError(http.StatusNotFound, 3003, "company_not_found")
}

// ------- Cluster Errors --------

func ClusterCreationFailed() (*Error) {
  return ApiError(http.StatusInternalServerError, 4001, "cluster_creation_failed")
}

func ClusterAlreadyExists() (*Error) {
  return ApiError(http.StatusInternalServerError, 4002, "cluster_already_exists")
}