package main

import (
  "fmt"
  "github.com/sweettea/rest-api/app"
  "github.com/sweettea/rest-api/pkg/utils"
  "github.com/sweettea/rest-api/pkg/database"
)

func main() {
  // Load app config
  utils.Assert(app.LoadConfig(), "Failed to load app config")

  // Establish connection to database
  db := database.Connection(app.Config.DatabaseUrl)

  fmt.Println(db)
}