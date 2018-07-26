package errmsg

import "net/http"

func ProjectConfigCreationFailed() (*Error) {
  return ApiError(http.StatusInternalServerError, 4001, "project_config_creation_failed")
}