package errmsg

import "net/http"

func ProjectNotAvailable() (*Error) {
  return ApiError(http.StatusInternalServerError, 5001, "project_not_available")
}

func ProjectCreationFailed() (*Error) {
  return ApiError(http.StatusInternalServerError, 5002, "project_creation_failed")
}