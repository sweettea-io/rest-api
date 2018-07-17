package testutil

import (
  "github.com/gorilla/mux"
  "net/http/httptest"
)

// Wrapper type around mux.Router, providing additional request functionality.
type TestRouter struct {
  Router *mux.Router
  BaseRoute string
  AuthHeaderName string
  AuthHeaderVal  string
}

// Perform an HTTP request, responding with a TestResponse object.
func (tr *TestRouter) Request(req *Request) *Response {
  // Create HTTP request object.
  httpReq := req.CreateHTTPRequest(tr.BaseRoute)

  if req.Authed {
    httpReq.Header.Set(tr.AuthHeaderName, tr.AuthHeaderVal)
  }

  // Create new response recorder.
  res := httptest.NewRecorder()

  // Perform the API call.
  tr.Router.ServeHTTP(res, httpReq)

  // Return a TestResponse wrapping the raw response.
  return &Response{raw: res}
}