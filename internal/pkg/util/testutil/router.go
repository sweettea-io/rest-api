package testutil

import (
  "github.com/gorilla/mux"
  "net/http/httptest"
)

// Wrapper type around mux.Router, providing additional request functionality.
type TestRouter struct {
  Router *mux.Router
}

// Perform an HTTP request, responding with a TestResponse object.
func (tr *TestRouter) Request(testReq *TestRequest) *TestResponse {
  // Create HTTP request object.
  req := testReq.CreateRequest()

  // Create new response recorder.
  res := httptest.NewRecorder()

  // Perform the API call.
  tr.Router.ServeHTTP(res, req)

  // Return a TestResponse wrapping the raw response.
  return &TestResponse{raw: res}
}