package route

import (
  "net/http"
  "github.com/sweettea-io/rest-api/internal/app/middleware"
)

// ----------- ROUTER SETUP ------------

// Prefix for all routes in this file.
const DeployRoute = "/deploy"

func InitDeployRouter() {
  // Create deploy router.
  deployRouter := Router.PathPrefix(DeployRoute).Subrouter()

  // Attach Session-based auth middleware to all request handlers on this router.
  deployRouter.Use(middleware.SessionAuth)

  // Attach route handlers.
  deployRouter.HandleFunc("", CreateModelVersionHandler).Methods("POST")
}

// ----------- ROUTE HANDLERS -----------

/*
  Create a Deploy

  Method:  POST
  Route:   /deploy
*/
func CreateDeployHandler(w http.ResponseWriter, req *http.Request) {
  // TODO: Handle file upload
}