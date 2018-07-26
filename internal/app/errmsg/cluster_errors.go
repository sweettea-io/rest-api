package errmsg

import "net/http"

func ClusterCreationFailed() (*Error) {
  return ApiError(http.StatusInternalServerError, 10001, "cluster_creation_failed")
}

func ClusterAlreadyExists() (*Error) {
  return ApiError(http.StatusInternalServerError, 10002, "cluster_already_exists")
}