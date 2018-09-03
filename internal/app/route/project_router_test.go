package route

import (
  "testing"
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/pkg/util/testutil"
  "github.com/sweettea-io/rest-api/internal/pkg/util/enc"
  "github.com/stretchr/testify/assert"
)

func TestUpsertProjectHandler(t *testing.T) {
  route := ProjectRoute

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
      Name: "invalid payload when unsupported project host used in namespace",
      Request: &testutil.Request{
        Method: "POST",
        Route: route,
        BeforeSend: []testutil.RequestModifier{
          testutil.AuthReqWithNewUser,
        },
        Data: &enc.JSON{
          "nsp": "unsupported-host.com/team/project",
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
      Name: "invalid payload when team section of namespace is not slugified",
      Request: &testutil.Request{
        Method: "POST",
        Route: route,
        BeforeSend: []testutil.RequestModifier{
          testutil.AuthReqWithNewUser,
        },
        Data: &enc.JSON{
          "nsp": "github.com/Team/project",
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
      Name: "invalid payload when repo section of namespace is not slugified",
      Request: &testutil.Request{
        Method: "POST",
        Route: route,
        BeforeSend: []testutil.RequestModifier{
          testutil.AuthReqWithNewUser,
        },
        Data: &enc.JSON{
          "nsp": "github.com/team/Project",
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
      Name: "does NOT create new project and fails when project already exists",
      Request: &testutil.Request{
        Method: "POST",
        Route: route,
        BeforeSend: []testutil.RequestModifier{
          testutil.AuthReqWithNewUser,
          testutil.CreateProject("github.com/team/project"),
        },
        Data: &enc.JSON{
          "nsp": "github.com/team/project",
        },
      },
      ExpectedStatus: 500,
      ExpectedRespJSON: &enc.JSON{
        "ok": false,
        "code": 3001,
        "error": "project_namespace_unavailable",
      },
      CustomAssertions: []testutil.CustomReqAssertion{
        func(t *testing.T, tc *testutil.RequestCase, status int, data *enc.JSON) {
          testutil.AssertTableCount(t, tc, "projects", 1)
        },
      },
    },
    {
      Name: "successfully creates project for proper, unique namespace",
      Request: &testutil.Request{
        Method: "POST",
        Route: route,
        BeforeSend: []testutil.RequestModifier{
          testutil.AuthReqWithNewUser,
        },
        Data: &enc.JSON{
          "nsp": "github.com/team/project",
        },
      },
      ExpectedStatus: 201,
      CustomAssertions: []testutil.CustomReqAssertion{
        func(t *testing.T, tc *testutil.RequestCase, status int, data *enc.JSON) {
          // Parse newly created Project.
          d := *data
          proj := d["project"].(map[string]interface{})
          projNsp := proj["nsp"].(string)

          // Assert new Project has proper slugified namespace.
          assert.Equal(t, "github.com/team/project", projNsp, tc.Name)

          // Ensure only one Project was created/exists with this namespace.
          testutil.AssertTableCount(t, tc, "projects", 1)
        },
      },
    },
  }

  EvalRequestCases(t, testCases)
}
