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
  // Initialize the app and start with a teardown.
  app.Init(config.New())
  Teardown()

  // Run the tests in this package and exit.
  code := m.Run()
  os.Exit(code)
}

// Setup test function for all tests in this package.
func Setup(cfg config.ConfigItf) {
  if cfg == nil {
    cfg = app.Config
  }

  InitRouter(cfg)

  // Create test router wrapping the base router.
  TestRouter = &testutil.Router{
    Raw: Router.GetRouter(),
    BaseRoute: cfg.BaseRoute(),
  }
}

// Teardown test function for all tests in this package.
func Teardown() {
  testutil.ClearTables(app.DB, true)
}
