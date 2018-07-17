package errmsg

import "net/http"

func SessionCreationFailed() (*Error) {
  return ApiError(http.StatusInternalServerError, 2001, "session_creation_failed")
}