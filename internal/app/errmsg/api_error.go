package errmsg

import (
  "github.com/sweettea-io/rest-api/internal/pkg/util/enc"
)

type Error struct {
  Status int
  Data enc.JSON
}

func ApiError(status int, code int, msg string) (*Error) {
  return &Error{
    Status: status,
    Data: enc.JSON{
      "ok": false,
      "code": code,
      "error": msg,
    },
  }
}

