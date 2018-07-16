package routes

import (
  "os"
  "testing"
  "github.com/Sirupsen/logrus"
  "github.com/sweettea-io/rest-api/app"
  "github.com/sweettea-io/rest-api/defs"
  "github.com/sweettea-io/rest-api/pkg/database"
  "github.com/sweettea-io/rest-api/pkg/test_utils"
)

var router *test_utils.TestRouter

func TestMain(m *testing.M) {
  // Load app config.
  app.LoadConfig()

  // Establish connection to database.
  db = database.Connection(app.Config.DatabaseUrl)
  db.LogMode(app.Config.Debug)

  // Clear all DB tables.
  test_utils.ClearTables(db, app.Config.Debug)

  // Create logger.
  logger = logrus.New()

  // Create the core router instance.
  apiRouter := CreateRouter(app.Config.BaseRoute(), db, logger)

  // Create API router.
  router = &test_utils.TestRouter{
    Router: apiRouter,
    BaseRoute: app.Config.BaseRoute(),
    AuthHeaderName: defs.AuthHeaderName,
    AuthHeaderVal: app.Config.RestApiToken,
  }

  // Run the tests in this package and exit.
  code := m.Run()
  os.Exit(code)
}

// Teardown function called after each test in this package is run.
// Call `defer Teardown()` at the top of each test you want to use this with.
func Teardown() {
  // Clear all DB tables.
  test_utils.ClearTables(db, app.Config.Debug)
}