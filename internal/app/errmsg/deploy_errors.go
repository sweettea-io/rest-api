package errmsg

import "net/http"

func DeployCreationFailed() (*Error) {
  return ApiError(http.StatusInternalServerError, 9001, "deploy_creation_failed")
}

func CreateDeploySchedulingFailed() (*Error) {
  return ApiError(http.StatusInternalServerError, 9002, "create_deploy_scheduling_failed")
}