package main

import (
  "fmt"
  "net/http"
  "github.com/Sirupsen/logrus"
  "github.com/sweettea/rest-api/app"
  "github.com/sweettea/rest-api/app/api"
  "github.com/sweettea/rest-api/pkg/database"
)

func main() {
  // Load app config.
  app.LoadConfig()

  // Create redis pool.
  app.CreateRedisPool()

  // Create job queue.
  app.CreateJobQueue()

  // Establish connection to database.
  db := database.Connection(app.Config.DatabaseUrl)
  db.LogMode(app.Config.Debug)

  // Create logger.
  logger := logrus.New()

  // Construct base route from API version (i.e. "/v1").
  baseRoute := fmt.Sprintf("/%s", app.Config.ApiVersion)

  // Create API router.
  router := api.CreateRouter(baseRoute, db, logger)

  // Format address to listen on.
  addr := fmt.Sprintf(":%v", app.Config.Port)

  logger.Infof("Listening on port %v...\n", app.Config.Port)

  // Start server.
  panic(http.ListenAndServe(addr, router))
}