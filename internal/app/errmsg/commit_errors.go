package errmsg

import "net/http"

func CommitCreationFailed() (*Error) {
  return ApiError(http.StatusInternalServerError, 5001, "commit_creation_failed")
}