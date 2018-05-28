package routes

import (
  "io"
  "net/http"
  "github.com/gorilla/mux"
  "github.com/sweettea/rest-api/pkg/utils/headers"
)

const UserRoute = "/users"

func InitUserRouter(baseRouter *mux.Router) {
  // Create user router
  userRouter := baseRouter.PathPrefix(UserRoute).Subrouter()

  // Attach route handlers
  userRouter.HandleFunc("", GetUsersHandler)
}

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusOK)
  w.Header().Set(headers.JsonContentType())
  io.WriteString(w, `{"alive": true}`)
}