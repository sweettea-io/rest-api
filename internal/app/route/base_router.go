package route

import (
  "github.com/gorilla/mux"
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/app/middleware"
)

var Router *mux.Router

func InitRouter() {
  // Create base router from provided baseRoute.
  Router = mux.NewRouter().PathPrefix(app.Config.BaseRoute()).Subrouter()

  // Attach base middleware.
  Router.Use(middleware.LogRequest)

  // Initialize sub routers.
  InitUserRouter()
  InitProjectRouter()
  InitTrainJobRouter()
  InitApiClusterRouter()
}