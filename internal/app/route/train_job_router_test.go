package route

import (
  "testing"
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/pkg/util/testutil"
  "github.com/sweettea-io/rest-api/internal/pkg/util/enc"
  "github.com/stretchr/testify/assert"
  "reflect"
  "github.com/sweettea-io/rest-api/internal/pkg/util/testutil/mocks"
)

func TestCreateTrainJobHandler(t *testing.T) {
  route := TrainJobRoute

  testCases := []testutil.RequestCase{
    {
      Name: "unauthorized when auth header not provided",
      Request: &testutil.Request{
        Method: "POST",
        Route: route,
      },
      ExpectedStatus: 401,
      ExpectedRespJSON: &enc.JSON{
        "ok": false,
        "code": 401,
        "error": "Unauthorized",

      },
    },
    {
      Name: "unauthorized when session token invalid",
      Request: &testutil.Request{
        Method: "POST",
        Route: route,
        Headers: map[string]string{
          app.Config.AuthHeaderName: "some-invalid-token",
        },
      },
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
        BeforeSend: []testutil.RequestModifier{
          testutil.AuthReqWithNewUser,
        },
      },
      ExpectedStatus: 400,
      ExpectedRespJSON: &enc.JSON{
        "ok": false,
        "code": 400,
        "error": "invalid_input_payload",
      },
    },
    {
      Name: "fails when train cluster not configured",
      CustomCfg: &mocks.MockConfig{
        MockTrainClusterConfigured: func() bool { return false },
      },
      Request: &testutil.Request{
        Method: "POST",
        Route: route,
        Data: &enc.JSON{
          "projectNsp": "my-project-nsp",
        },
        BeforeSend: []testutil.RequestModifier{
          testutil.AuthReqWithNewUser,
        },
      },
      ExpectedStatus: 500,
      ExpectedRespJSON: &enc.JSON{
        "ok": false,
        "code": 6002,
        "error": "train_cluster_not_configured",
      },
    },
  }

  for _, tc := range testCases {
    func () {
      Setup(tc.SetupArgs())
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