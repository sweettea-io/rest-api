package main

import (
  "fmt"
  "net/http"
  "github.com/Sirupsen/logrus"
  "github.com/sweettea-io/rest-api/app"
  "github.com/sweettea-io/rest-api/pkg/database"
  "github.com/sweettea-io/rest-api/app/api/routes"
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

  // Create API router.
  router := routes.CreateRouter(app.Config.BaseRoute(), db, logger)

  // Format address to listen on.
  addr := fmt.Sprintf(":%v", app.Config.ServerPort)

  logger.Infof("Listening on port %v...\n", app.Config.ServerPort)

  // Start server.
  panic(http.ListenAndServe(addr, router))
}