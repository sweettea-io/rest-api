package main

import (
  "github.com/sweettea/rest-api/app"
  "github.com/sweettea/rest-api/pkg/database"
)

func main() {
  // Load app config
  app.LoadConfig()

  // Migrate any DB changes
  database.Migrate(app.Config.DatabaseUrl)
}