package route

import (
  "testing"
  "os"
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/pkg/config"
  "github.com/sweettea-io/rest-api/internal/pkg/util/testutil"
  "github.com/stretchr/testify/assert"
  "reflect"
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

func EvalRequestCases(t *testing.T, testCases []testutil.RequestCase) {
  for _, tc := range testCases {
    func () {
      // Setup and defer Teardown.
      Setup(tc.SetupArgs())
      defer Teardown()

      // Perform request and get response.
      res, err := TestRouter.Request(tc.Request)
      if err != nil {
        t.Error(err.Error())
        return
      }

      // Get response status and parse data.
      status := res.StatusCode()
      data := res.ParseJSON()

      // Assert request status if provided.
      if tc.ExpectedStatus != 0 {
        assert.Equal(t, tc.ExpectedStatus, status, tc.Name)
      }

      // Assert response JSON equality if provided.
      if tc.ExpectedRespJSON != nil {
        assert.Equal(t, true, reflect.DeepEqual(tc.ExpectedRespJSON.Cycle(), data), tc.Name)
      }

      // Perform custom assertions if provided.
      if tc.CustomAssertions != nil && len(tc.CustomAssertions) > 0 {
        for _, assertion := range tc.CustomAssertions {
          assertion(t, &tc, status, data)
        }
      }
    }()
  }
}