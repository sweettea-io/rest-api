package error

import "net/http"

type Error map[string]interface{}

func Forbidden() (int, *Error) {
  return http.StatusForbidden, &Error{"ok": false, "code": http.StatusForbidden, "error": "forbidden"}
}