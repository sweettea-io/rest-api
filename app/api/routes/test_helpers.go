//
// This file includes helper types and functions for testing API routes.
//
package routes

import (
  "encoding/json"
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

// --------------- TEST ROUTER ---------------

// Wrapper type around mux.Router, providing additional request functionality.
type testRouter struct {
  Router *mux.Router
}

// Initialize the Router property with a new mux.Router.
func (tr *testRouter) Init() {
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

// Perform an HTTP request, responding with a testResponse object.
func (tr *testRouter) Request(method string, route string, body io.Reader, authed bool, extraHeaders ...map[string]string) *testResponse {
  // Init the mux router if non-existent.
  if tr.Router == nil {
    tr.Init()
  }

  // Get extra headers from arg (if they exist).
  var headers map[string]string

  if len(extraHeaders) > 0 {
    headers = extraHeaders[0]
  }

  // Create HTTP request object.
  req := tr.createRequest(method, route, body, authed, headers)

  // Create new response recorder.
  res := httptest.NewRecorder()

  // Perform the API call.
  tr.Router.ServeHTTP(res, req)

  // Return a testResponse wrapping the raw response.
  return &testResponse{raw: res}
}

func (tr *testRouter) JSONRequest(method string, route string, data *utils.JSON, authed bool, extraHeaders ...map[string]string) *testResponse {
  // Init the mux router if non-existent.
  if tr.Router == nil {
    tr.Init()
  }

  var body io.Reader = nil

  // Convert JSON data into a buffer if it exists to represent the request body.
  if data != nil {
    body, _ = data.AsBuffer()
  }

  // Get extra headers from arg (if they exist).
  var headers map[string]string

  if len(extraHeaders) > 0 {
    headers = extraHeaders[0]
  }

  // Create HTTP request object.
  req := tr.createRequest(method, route, body, authed, headers)

  // Configure request body to be of JSON type.
  req.Header.Set("Content-Type", "application/json")

  // Create new response recorder.
  res := httptest.NewRecorder()

  // Perform the API call.
  tr.Router.ServeHTTP(res, req)

  // Return a testResponse wrapping the raw response.
  return &testResponse{raw: res}
}

// Create and return a new HTTP request object for the test router.
func (tr *testRouter) createRequest(method string, route string, body io.Reader, authed bool, extraHeaders map[string]string) *http.Request {
  // Create new HTTP request.
  req, err := http.NewRequest(method, app.Config.BaseRoute() + route, body)

  if err != nil {
    panic(fmt.Errorf("error creating new http request object: %s", err.Error()))
  }

  // Add Sweet Tea auth header if specified as an 'authed' request.
  if authed {
    req.Header.Set(defs.AuthHeaderName, app.Config.RestApiToken)
  }

  // Set any extra headers provided.
  for k, v := range extraHeaders {
    req.Header.Set(k, v)
  }

  return req
}

// Router pointer used by all *_test.go router test files in this package.
var TestRouter = &testRouter{}

// --------------- TEST RESPONSE ---------------

// Wrapper type around httptest.ResponseRecorder, providing additional helper functions for
// JSON parsing the response body and fetching the response status code.
type testResponse struct {
  raw *httptest.ResponseRecorder
}

// Parse and return the response body as utils.JSON type.
func (res *testResponse) ParseJSON() *utils.JSON {
  // Parse body into JSON type.
  var data utils.JSON
  json.Unmarshal(res.raw.Body.Bytes(), &data)

  return &data
}

// Parse and return the response status code.
func (res *testResponse) StatusCode() int {
  return res.raw.Result().StatusCode
}