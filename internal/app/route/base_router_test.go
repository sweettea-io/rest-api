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

  // Setup test router.
  InitRouter()
  TestRouter = &testutil.Router{Raw: Router}

  // Run the tests in this package and exit.
  code := m.Run()
  os.Exit(code)
}

// Teardown test function for all tests in this package.
func Teardown() {
  testutil.ClearTables(app.DB, app.Config.Debug)
}
