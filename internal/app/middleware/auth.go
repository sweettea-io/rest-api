package middleware

import (
  "net/http"
  "github.com/sweettea-io/rest-api/internal/app/errmsg"
  "github.com/sweettea-io/rest-api/internal/app/respond"
  "github.com/sweettea-io/rest-api/internal/pkg/service/usersvc"
)

func SessionAuth(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
    // Auth request from Session token.
    _, err := usersvc.FromRequest(req)

    if err != nil {
      respond.Error(w, errmsg.Unauthorized())
      return
    }

    // Call the next handler.
    next.ServeHTTP(w, req)
  })
}