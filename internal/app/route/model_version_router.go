package route

import (
  "net/http"
  "github.com/sweettea-io/rest-api/internal/app/middleware"
)

// ----------- ROUTER SETUP ------------

// Prefix for all routes in this file.
const ModelVersionRouter = "/model_version"

func InitModelVersionRouter() {
  // Create modelVersion router.
  modelVersion := Router.PathPrefix(ModelVersionRouter).Subrouter()

  // Attach Session-based auth middleware to all request handlers on this router.
  modelVersion.Use(middleware.SessionAuth)

  // Attach route handlers.
  modelVersion.HandleFunc("", CreateModelVersionHandler).Methods("POST")
}

// ----------- ROUTE HANDLERS -----------

/*
  Create a ModelVersion

  Method:  POST
  Route:   /model_version
*/
func CreateModelVersionHandler(w http.ResponseWriter, req *http.Request) {
  // TODO: Handle file upload
}