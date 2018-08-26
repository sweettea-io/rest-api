package route

import (
  "testing"
  "os"
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/pkg/config"
  "github.com/sweettea-io/rest-api/internal/pkg/util/testutil"
)

var TestRouter *testutil.Router

func TestMain(m *testing.M) {
  // Initialize the app, start with a teardown, and create the test router.
  app.Init(config.New())
  Teardown()
  InitRouter()

  // Create test router wrapping the base router.
  TestRouter = &testutil.Router{
    Raw: Router,
    BaseRoute: app.Config.BaseRoute(),
  }

  // Run the tests in this package and exit.
  code := m.Run()
  os.Exit(code)
}

// Teardown test function for all tests in this package.
func Teardown() {
  testutil.ClearTables(app.DB, app.Config.Debug)
}
