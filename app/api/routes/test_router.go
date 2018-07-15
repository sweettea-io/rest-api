package routes

import (
  "fmt"
  "io"
  "net/http"
  "net/http/httptest"
  "github.com/gorilla/mux"
  "github.com/Sirupsen/logrus"
  "github.com/sweettea-io/rest-api/app"
  "github.com/sweettea-io/rest-api/defs"
  "github.com/sweettea-io/rest-api/pkg/database"
  "github.com/sweettea-io/rest-api/pkg/utils"
)

type testRouter struct {
  Router *mux.Router
}

func (tr *testRouter) Configure() {
  // Load app config.
  app.LoadConfig()

  // Establish connection to database.
  db := database.Connection(app.Config.DatabaseUrl)
  db.LogMode(app.Config.Debug)

  // Create logger.
  logger := logrus.New()

  // Create API router.
  tr.Router = CreateRouter(app.Config.BaseRoute(), db, logger)
}

func (tr *testRouter) Request(method string, route string, body io.Reader, authed bool) *httptest.ResponseRecorder {
  if tr.Router == nil {
    tr.Configure()
  }

  req := newRequestObj(method, route, body, authed)

  res := httptest.NewRecorder()

  tr.Router.ServeHTTP(res, req)

  return res
}

func (tr *testRouter) JSONRequest(method string, route string, data *utils.JSON, authed bool) *httptest.ResponseRecorder {
  if tr.Router == nil {
    tr.Configure()
  }

  var body io.Reader = nil

  if data != nil {
    body, _ = data.AsBuffer()
  }

  req := newRequestObj(method, route, body, authed)

  req.Header.Set("Content-Type", "application/json")

  res := httptest.NewRecorder()

  tr.Router.ServeHTTP(res, req)

  // TODO:

  return res
}

func newRequestObj(method string, route string, body io.Reader, authed bool) *http.Request {
  req, err := http.NewRequest(method, app.Config.BaseRoute() + route, body)

  if err != nil {
    panic(fmt.Errorf("error creating new http request object: %s", err.Error()))
  }

  if authed {
    req.Header.Set(defs.AuthHeaderName, app.Config.RestApiToken)
  }

  return req
}

var TestRouter = &testRouter{}