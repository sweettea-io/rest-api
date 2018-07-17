package routes

import (
  "os"
  "testing"
  "github.com/Sirupsen/logrus"
  "github.com/sweettea-io/rest-api/app"
  "github.com/sweettea-io/rest-api/pkg/database"
  "github.com/sweettea-io/rest-api/pkg/test_utils"
)

// Global test router used by tests as wrapper to *mux.Router
var router *test_utils.TestRouter

// Setup function called before all tests in this package are run.
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

  // Create API router.
  router = &test_utils.TestRouter{
    Router: CreateRouter(app.Config.BaseRoute(), db, logger),
  }

  // Run the tests in this package and exit.
  code := m.Run()
  os.Exit(code)
}

// RouteTeardown function called after each test in this package is run.
// Call `defer RouteTeardown()` at the top of each test you want to use this with.
func RouteTeardown() {
  // Clear all DB tables.
  test_utils.ClearTables(db, app.Config.Debug)
}