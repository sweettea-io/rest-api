package errmsg

import "net/http"

func CompanyAlreadyExists() (*Error) {
  return ApiError(http.StatusInternalServerError, 3001, "company_already_exists")
}

func CompanyCreationFailed() (*Error) {
  return ApiError(http.StatusInternalServerError, 3002, "company_creation_failed")
}

func CompanyNotFound() (*Error) {
  return ApiError(http.StatusNotFound, 3003, "company_not_found")
}
