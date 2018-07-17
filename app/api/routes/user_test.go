package routes

import (
  "testing"
  "github.com/sweettea-io/rest-api/pkg/test_utils"
  "github.com/sweettea-io/rest-api/pkg/test_utils/cases"
  "github.com/sweettea-io/rest-api/pkg/utils"
  "github.com/stretchr/testify/assert"
  "reflect"
)

func TestCreateUserHandler(t *testing.T) {
  // Create list of JSON request test cases.
  testCases := []cases.RequestCase{
    {
      Name: "request unauthorized when auth header NOT provided",
      Request: &test_utils.TestRequest{
        Method: "POST",
        Route: UserRoute,
      },
      ExpectedStatus: 401,
      ExpectedRespJSON: &utils.JSON{
        "ok": false,
        "code": 401,
        "error": "Unauthorized",
      },
    },
  }

  // Evaluate each test case and teardown after each.
  for _, tc  := range testCases {
    func () {
      defer RouteTeardown()

      // Make request.
      res := router.Request(tc.Request)

      // Assert response code and data.
      assert.Equal(t, tc.ExpectedStatus, res.StatusCode(), tc.Name)
      assert.Equal(t, true, reflect.DeepEqual(tc.ExpectedRespJSON.Cycle(), res.ParseJSON()), tc.Name)
    }()
  }
}