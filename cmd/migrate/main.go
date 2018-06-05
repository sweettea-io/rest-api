package main

import (
  "github.com/sweettea-io/rest-api/app"
  "github.com/sweettea-io/rest-api/pkg/database"
)

func main() {
  // Load app config
  app.LoadConfig()

  // Migrate any DB changes
  database.Migrate(app.Config.DatabaseUrl)
}