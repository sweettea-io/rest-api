package route

import (
  "github.com/gorilla/mux"
  "github.com/sweettea-io/rest-api/internal/app/middleware"
  "github.com/sweettea-io/rest-api/internal/pkg/config"
)

type ConfiguredRouter struct {
  Internal *mux.Router
  Config config.ConfigItf
}

func (cr *ConfiguredRouter) GetRouter() *mux.Router {
  return cr.Internal
}

var Router *ConfiguredRouter

func InitRouter(cfg config.ConfigItf) {
  // Create base router from provided baseRoute.
  Router = &ConfiguredRouter{
    Internal: mux.NewRouter().PathPrefix(cfg.BaseRoute()).Subrouter(),
    Config: cfg,
  }

  // Attach base middleware.
  Router.GetRouter().Use(middleware.LogRequest)

  // Initialize sub routers.
  InitUserRouter()
  InitProjectRouter()
  InitTrainJobRouter()
  InitDeployRouter()
  InitModelVersionRouter()
  InitApiClusterRouter()
}

