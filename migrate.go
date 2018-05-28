package main

import (
  "github.com/sweettea/rest-api/app"
  "github.com/sweettea/rest-api/pkg/database"
  "github.com/sweettea/rest-api/pkg/utils"
)

func main() {
  // Load app config
  utils.Assert(app.LoadConfig(), "Failed to load app config")

  // Migrate any DB changes
  database.Migrate(app.Config.DatabaseUrl)
}