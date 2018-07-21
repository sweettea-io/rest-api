package route

import (
  "reflect"
  "testing"
  "github.com/stretchr/testify/assert"
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
      Name: "bad req data results in invalid payload",
      Request: &testutil.Request{
        Method: "POST",
        Route: route,
        Data: &enc.JSON{},
        Authed: true,
      },
      ExpectedStatus: 400,
      ExpectedRespJSON: &enc.JSON{"ok": false, "code": 400, "error": "invalid_input_payload"},
    },
  }

  // TODO: For the 2nd case above, what you would want to do is stub CreateUserPayload.Validate() to return false, ...
  // TODO: ...and then check that what you have is returned.

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