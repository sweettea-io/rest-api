package testutil

import (
  "github.com/gorilla/mux"
  "net/http/httptest"
)

// Wrapper type around mux.Router, providing additional request functionality.
type Router struct {
  Raw *mux.Router
  BaseRoute string
  AuthHeaderName string
  AuthHeaderVal  string
}

// Perform an HTTP request, responding with a Response object.
func (router *Router) Request(req *Request) *Response {
  // Create HTTP request object.
  httpReq := req.CreateHTTPRequest(router.BaseRoute)

  if req.Authed {
    httpReq.Header.Set(router.AuthHeaderName, router.AuthHeaderVal)
  }

  // Create new response recorder.
  res := httptest.NewRecorder()

  // Perform the API call.
  router.Raw.ServeHTTP(res, httpReq)

  // Return a Response wrapping the raw response.
  return &Response{raw: res}
}