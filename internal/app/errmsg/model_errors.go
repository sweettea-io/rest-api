package errmsg

import "net/http"

func ModelCreationFailed() (*Error) {
  return ApiError(http.StatusInternalServerError, 7001, "model_creation_failed")
}