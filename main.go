package main

import (
  "fmt"
  "net/http"
  "github.com/sweettea/rest-api/app"
  "github.com/sweettea/rest-api/app/api"
  "github.com/sweettea/rest-api/pkg/database"
  "github.com/sweettea/rest-api/pkg/utils"
  "github.com/Sirupsen/logrus"
)

func main() {
  // Load app config.
  utils.Assert(app.LoadConfig(), "Failed to load app config")

  // Establish connection to database.
  db := database.Connection(app.Config.DatabaseUrl)

  // Create logger
  logger := logrus.New()

  // Construct the base route from the API version (i.e. "/v1").
  baseRoute := fmt.Sprintf("/%s", app.Config.ApiVersion)

  // Create API router.
  router := api.CreateRouter(baseRoute, db, logger)

  // Format address to listen on.
  addr := fmt.Sprintf(":%v", app.Config.Port)

  logger.Infof("Listening on port %v...\n", app.Config.Port)

  // Start the server.
  panic(http.ListenAndServe(addr, router))
}