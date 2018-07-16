package test_utils

import (
  "github.com/gorilla/mux"
  "github.com/sweettea-io/rest-api/pkg/utils"
  "io"
  "net/http/httptest"
  "net/http"
  "fmt"
)

// Wrapper type around mux.Router, providing additional request functionality.
type TestRouter struct {
  Router         *mux.Router
  BaseRoute      string
  AuthHeaderName string
  AuthHeaderVal  string
}

// Perform an HTTP request, responding with a TestResponse object.
func (tr *TestRouter) Request(method string, route string, body io.Reader, authed bool, extraHeaders ...map[string]string) *TestResponse {
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

  // Return a TestResponse wrapping the raw response.
  return &TestResponse{raw: res}
}

func (tr *TestRouter) JSONRequest(method string, route string, data *utils.JSON, authed bool, extraHeaders ...map[string]string) *TestResponse {
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

  // Return a TestResponse wrapping the raw response.
  return &TestResponse{raw: res}
}

// Create and return a new HTTP request object for the test router.
func (tr *TestRouter) createRequest(method string, route string, body io.Reader, authed bool, extraHeaders map[string]string) *http.Request {
  // Create new HTTP request.
  req, err := http.NewRequest(method, tr.BaseRoute + route, body)

  if err != nil {
    panic(fmt.Errorf("error creating new http request object: %s", err.Error()))
  }

  // Add Sweet Tea auth header if specified as an 'authed' request.
  if authed {
    req.Header.Set(tr.AuthHeaderName, tr.AuthHeaderVal)
  }

  // Set any extra headers provided.
  for k, v := range extraHeaders {
    req.Header.Set(k, v)
  }

  return req
}