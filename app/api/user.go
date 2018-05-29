package api

import (
  "net/http"
  "github.com/gorilla/mux"
  "github.com/sweettea/rest-api/app/api/resp"
)

const UserRoute = "/users"

func InitUserRouter(baseRouter *mux.Router) {
  // Create user router
  userRouter := baseRouter.PathPrefix(UserRoute).Subrouter()

  // Attach route handlers
  userRouter.HandleFunc("/auth", UserAuthHandler).Methods("GET")
}

func UserAuthHandler(w http.ResponseWriter, req *http.Request) {
  respOk(w, resp.UserLoginSuccess)
}