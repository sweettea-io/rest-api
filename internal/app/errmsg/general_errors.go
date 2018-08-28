package errmsg

import "net/http"

const JsonEncodingError = `{"ok": false, "code": 1000, "error": "Error encoding JSON response"}`

func Forbidden() (*Error) {
  status := http.StatusForbidden
  return ApiError(status, status, http.StatusText(status))
}

func Unauthorized() (*Error) {
  status := http.StatusUnauthorized
  return ApiError(status, status, http.StatusText(status))
}

func ISE() (*Error) {
  status := http.StatusInternalServerError
  return ApiError(status, status, http.StatusText(status))
}

func InvalidPayload() (*Error) {
  return ApiError(http.StatusBadRequest, http.StatusBadRequest, "invalid_input_payload")
}

func StreamingNotSupported() (*Error) {
  return ApiError(http.StatusBadRequest, 999, "streaming_not_supported")
}