package errmsg

import "net/http"

func ClusterCreationFailed() (*Error) {
  return ApiError(http.StatusInternalServerError, 4001, "cluster_creation_failed")
}

func ClusterAlreadyExists() (*Error) {
  return ApiError(http.StatusInternalServerError, 4002, "cluster_already_exists")
}