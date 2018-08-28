package route

import (
  "reflect"
  "testing"
  "github.com/stretchr/testify/assert"
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/pkg/util/testutil"
  "github.com/sweettea-io/rest-api/internal/pkg/util/enc"
)

func TestCreateUserHandler(t *testing.T) {
  route := UserRoute

  testCases := []testutil.RequestCase{
    {
      Name: "unauthorized when auth header NOT provided",
      Request: &testutil.Request{Method: "POST", Route: route},
      ExpectedStatus: 401,
      ExpectedRespJSON: &enc.JSON{"ok": false, "code": 401, "error": "Unauthorized"},
    },
    {
      Name: "invalid payload when required fields missing",
      Request: &testutil.Request{
        Method: "POST",
        Route: route,
        Headers: map[string]string{
          app.Config.AuthHeaderName: app.Config.RestApiToken,
        },
      },
      ExpectedStatus: 400,
      ExpectedRespJSON: &enc.JSON{"ok": false, "code": 400, "error": "invalid_input_payload"},
    },
  }

  for _, tc := range testCases {
    func () {
      defer Teardown()

      // Perform request.
      res, err := TestRouter.Request(tc.Request)
      if err != nil {
        t.Error(err.Error())
        return
      }

      // Assert response status code and data.
      assert.Equal(t, tc.ExpectedStatus, res.StatusCode(), tc.Name)
      assert.Equal(t, true, reflect.DeepEqual(tc.ExpectedRespJSON.Cycle(), res.ParseJSON()), tc.Name)
    }()
  }
}