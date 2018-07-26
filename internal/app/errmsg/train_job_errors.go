package errmsg

import "net/http"

func TrainJobCreationFailed() (*Error) {
  return ApiError(http.StatusInternalServerError, 6001, "train_job_creation_failed")
}

func TrainClusterNotConfigured() (*Error) {
  return ApiError(http.StatusInternalServerError, 6002, "train_cluster_not_configured")
}

func CreateTrainJobSchedulingFailed() (*Error) {
  return ApiError(http.StatusInternalServerError, 6003, "create_train_job_scheduling_failed")
}