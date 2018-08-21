package main

import (
  "fmt"
  "net/http"
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/app/route"
  "github.com/sweettea-io/rest-api/internal/pkg/config"
  "github.com/sweettea-io/rest-api/internal/pkg/util/cloud"
)

func main() {
  // Initialize the app.
  app.Init(config.New())

  // Create cloud session object.
  cloud.InitCloud()

  // Build API routes.
  route.InitRouter()

  // Format address to listen on.
  addr := fmt.Sprintf(":%v", app.Config.ServerPort)

  app.Log.Infof("Listening on port %v...\n", app.Config.ServerPort)

  // Start server.
  panic(http.ListenAndServe(addr, route.Router))
}