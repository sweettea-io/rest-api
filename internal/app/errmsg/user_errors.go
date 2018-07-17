package errmsg

import "net/http"

func UserNotFound() (*Error) {
  return ApiError(http.StatusNotFound, 1001, "user_not_found")
}

func EmailNotAvailable() (*Error) {
  return ApiError(http.StatusInternalServerError, 1002, "email_not_available")
}

func UserCreationFailed() (*Error) {
  return ApiError(http.StatusInternalServerError, 1003, "user_creation_failed")
}