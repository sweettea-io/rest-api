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
  // Initialize the app.
  app.Init(config.New())

  // Start with a teardown.
  Teardown()

  // Initialize base mux Router.
  InitRouter()

  // Create test router wrapping the base router.
  TestRouter = &testutil.Router{
    Raw: Router,
    BaseRoute: app.Config.BaseRoute(),
    AuthHeaderName: app.Config.AuthHeaderName,
    AuthHeaderVal: app.Config.RestApiToken,
  }

  // Run the tests in this package and exit.
  code := m.Run()
  os.Exit(code)
}

// Teardown test function for all tests in this package.
func Teardown() {
  testutil.ClearTables(app.DB, app.Config.Debug)
}
