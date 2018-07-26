package errmsg

import "net/http"

func EnvVarCreationFailed() (*Error) {
  return ApiError(http.StatusInternalServerError, 11001, "env_var_creation_failed")
}