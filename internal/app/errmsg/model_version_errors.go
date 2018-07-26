package errmsg

import "net/http"

func ModelVersionCreationFailed() (*Error) {
  return ApiError(http.StatusInternalServerError, 8001, "model_version_creation_failed")
}