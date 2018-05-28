package routes

import (
  "io"
  "net/http"
  "github.com/gorilla/mux"
)

const UserRoute = "/users"

func InitUserRouter(baseRouter *mux.Router) {
  // Create user router
  userRouter := baseRouter.PathPrefix(UserRoute).Subrouter()

  // Attach route handlers
  userRouter.HandleFunc("/", GetUsersHandler)
}

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
  // A very simple health check.
  w.WriteHeader(http.StatusOK)
  w.Header().Set("Content-Type", "application/json")

  // In the future we could report back on the status of our DB, or our cache
  // (e.g. Redis) by performing a simple PING, and include them in the response.
  io.WriteString(w, `{"alive": true}`)
}