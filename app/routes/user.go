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

func GetUsersHandler(writer http.ResponseWriter, req *http.Request) {
  writer.WriteHeader(http.StatusOK)
  writer.Header().Set(headers.JsonContentType())
  io.WriteString(writer, `{"alive": true}`)
}