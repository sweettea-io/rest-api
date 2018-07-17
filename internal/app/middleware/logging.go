package middleware

import (
  "net/http"
  "github.com/sweettea-io/rest-api/internal/app"
)

// LogRequest logs an HTTP request and continues to the next handler.
func LogRequest(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
    app.Log.Infoln(req.RequestURI)

    // Call the next handler.
    next.ServeHTTP(w, req)
  })
}