package route

import (
  "testing"
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
      ExpectedRespJSON: &enc.JSON{
        "ok": false,
        "code": 401,
        "error": "Unauthorized",
      },
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
      ExpectedRespJSON: &enc.JSON{
        "ok": false,
        "code": 400,
        "error": "invalid_input_payload",
      },
    },
  }

  EvalRequestCases(t, testCases)
}