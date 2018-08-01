package errmsg

import "net/http"

func ModelCreationFailed() (*Error) {
  return ApiError(http.StatusInternalServerError, 7001, "model_creation_failed")
}

func ModelNotFound() (*Error) {
  return ApiError(http.StatusNotFound, 7002, "model_not_found")
}

