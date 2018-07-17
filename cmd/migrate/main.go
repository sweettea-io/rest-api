package main

import (
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/pkg/config"
  "github.com/sweettea-io/rest-api/internal/pkg/db"
)

func main() {
  // Initialize the app.
  app.Init(config.New())

  // Auto-migrate any DB changes.
  db.Migrate(app.DB)
}