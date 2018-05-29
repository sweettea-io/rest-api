package api

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
  userRouter.HandleFunc("", GetUsersHandler)
}

func GetUsersHandler(w http.ResponseWriter, req *http.Request) {
  w.WriteHeader(http.StatusOK)
  io.WriteString(w, `{"alive": true}`)
}