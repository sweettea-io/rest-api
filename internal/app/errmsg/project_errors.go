package errmsg

import "net/http"

func ProjectNspUnavailable() (*Error) {
  return ApiError(http.StatusInternalServerError, 3001, "project_namespace_unavailable")
}

func ProjectUpsertionFailed() (*Error) {
  return ApiError(http.StatusInternalServerError, 3002, "project_upsertion_failed")
}