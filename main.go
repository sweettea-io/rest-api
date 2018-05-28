package main

import (
  "fmt"
  "net/http"
  "github.com/sweettea/rest-api/app"
  "github.com/sweettea/rest-api/app/routes"
  "github.com/sweettea/rest-api/pkg/database"
  "github.com/sweettea/rest-api/pkg/utils"
  logger "github.com/Sirupsen/logrus"
)

func main() {
  // Load app config
  utils.Assert(app.LoadConfig(), "Failed to load app config")

  // Establish connection to database
  db := database.Connection(app.Config.DatabaseUrl)

  router := routes.InitBaseRouter(db)

  addr := fmt.Sprintf(":%v", app.Config.Port)
  logger.Infof("Listening on port %v...\n", app.Config.Port)

  panic(http.ListenAndServe(addr, router))
}