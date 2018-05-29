package api

import (
  "net/http"
)

// Log the request and continue.
func LogRequest(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
    logger.Infoln(req.RequestURI)

    // Call the next handler.
    next.ServeHTTP(w, req)
  })
}