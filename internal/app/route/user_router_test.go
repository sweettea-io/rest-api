package route

import (
  "reflect"
  "testing"
  "github.com/stretchr/testify/assert"
  "github.com/sweettea-io/rest-api/internal/pkg/util/testutil"
  "github.com/sweettea-io/rest-api/internal/pkg/util/enc"
)

func TestCreateUserHandler(t *testing.T) {
  testCases := []testutil.RequestCase{
    {
      Name: "request unauthorized when auth header NOT provided",
      Request: &testutil.Request{Method: "POST", Route: UserRoute},
      ExpectedStatus: 401,
      ExpectedRespJSON: &enc.JSON{"ok": false, "code": 401, "error": "Unauthorized"},
    },
  }

  for _, tc := range testCases {
    func () {
      defer Teardown()

      // Perform request.
      res := TestRouter.Request(tc.Request)

      // Assert response status code and data.
      assert.Equal(t, tc.ExpectedStatus, res.StatusCode(), tc.Name)
      assert.Equal(t, true, reflect.DeepEqual(tc.ExpectedRespJSON.Cycle(), res.ParseJSON()), tc.Name)
    }()
  }
}