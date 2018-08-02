package errmsg

import "net/http"

func ModelVersionCreationFailed() (*Error) {
  return ApiError(http.StatusInternalServerError, 8001, "model_version_creation_failed")
}

func ModelVersionNotFound() (*Error) {
  return ApiError(http.StatusNotFound, 8002, "model_version_not_found")
}