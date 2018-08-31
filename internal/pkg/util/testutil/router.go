package testutil

import (
  "github.com/gorilla/mux"
  "net/http/httptest"
)

// Wrapper type around mux.Router, providing additional request functionality.
type Router struct {
  Raw *mux.Router
  BaseRoute string
}

// Perform an HTTP request, responding with a Response object.
func (router *Router) Request(req *Request) (*Response, error) {
  // Call before-send handler if exists.
  if req.BeforeSend != nil && len(req.BeforeSend) > 0 {
    var err error

    for _, bs := range req.BeforeSend {
      req, err = bs(req)

      if err != nil {
        return nil, err
      }
    }
  }

  // Create HTTP request object.
  httpReq := req.CreateHTTPRequest(router.BaseRoute)

  // Create new response recorder.
  res := httptest.NewRecorder()

  // Perform the API call.
  router.Raw.ServeHTTP(res, httpReq)

  // Return a Response wrapping the raw response.
  return &Response{raw: res}, nil
}