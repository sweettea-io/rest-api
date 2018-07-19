package route

import (
  "net/http"
  "github.com/sweettea-io/rest-api/internal/pkg/service/usersvc"
  "github.com/sweettea-io/rest-api/internal/app/respond"
  "github.com/sweettea-io/rest-api/internal/app/errmsg"
)

// ----------- ROUTER SETUP ------------

// Prefix for all routes in this file
const ProjectRoute = "/projects"

func InitProjectRouter() {
  // Create project router.
  projectRouter := Router.PathPrefix(ProjectRoute).Subrouter()

  // Attach route handlers.
  projectRouter.HandleFunc("", CreateProjectHandler).Methods("POST")
  projectRouter.HandleFunc("", ReadProjectHandler).Methods("GET")
  projectRouter.HandleFunc("", UpdateProjectHandler).Methods("UPDATE")
  projectRouter.HandleFunc("", DeleteProjectHandler).Methods("DELETE")
}

func CreateProjectHandler(w http.ResponseWriter, req *http.Request) {
  // Find user from request session.
  user, err := usersvc.FromRequest(req)

  if err != nil {
    respond.Error(w, errmsg.UserNotFound())
  }

  // Do user shit
}

func ReadProjectHandler(w http.ResponseWriter, req *http.Request) {
  // Find user from request session.
  user, err := usersvc.FromRequest(req)

  if err != nil {
    respond.Error(w, errmsg.UserNotFound())
  }

  // Do user shit
}

func UpdateProjectHandler(w http.ResponseWriter, req *http.Request) {
  // Find user from request session.
  user, err := usersvc.FromRequest(req)

  if err != nil {
    respond.Error(w, errmsg.UserNotFound())
  }

  // Do user shit
}

func DeleteProjectHandler(w http.ResponseWriter, req *http.Request) {
  // Find user from request session.
  user, err := usersvc.FromRequest(req)

  if err != nil {
    respond.Error(w, errmsg.UserNotFound())
  }

  // Do user shit
}