package errmsg

import "net/http"

func ApiClusterCreationFailed() (*Error) {
  return ApiError(http.StatusInternalServerError, 10001, "api_cluster_creation_failed")
}

func ApiClusterAlreadyExists() (*Error) {
  return ApiError(http.StatusInternalServerError, 10002, "api_cluster_already_exists")
}

func ApiClusterNotFound() (*Error) {
  return ApiError(http.StatusNotFound, 10003, "api_cluster_not_found")
}

func ApiClusterUpdateFailed() (*Error) {
  return ApiError(http.StatusInternalServerError, 10004, "api_cluster_update_failed")
}

func ApiClusterDeletionFailed() (*Error) {
  return ApiError(http.StatusInternalServerError, 10005, "api_cluster_deletion_failed")
}