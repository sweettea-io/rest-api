package errmsg

import "net/http"

func TrainJobCreationFailed() (*Error) {
  return ApiError(http.StatusInternalServerError, 6001, "train_job_creation_failed")
}