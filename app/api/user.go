package api

import (
  "net/http"
  "github.com/gorilla/mux"
)

const UserRoute = "/user"

func InitUserRouter(baseRouter *mux.Router) {
  // Create user router
  userRouter := baseRouter.PathPrefix(UserRoute).Subrouter()

  // Attach route handlers
  userRouter.HandleFunc("/auth", UserAuthHandler)
}

func UserAuthHandler(w http.ResponseWriter, req *http.Request) {
  user := CurrentUser(w, req)
  logger.Info(user)
}