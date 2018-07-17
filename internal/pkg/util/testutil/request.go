package testutil

import (
  "io"
  "net/http"
  "fmt"
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/pkg/util/enc"
)

type TestRequest struct {
  Method       string
  Route        string
  Body         io.Reader
  Data         *enc.JSON
  Authed       bool
  ExtraHeaders map[string]string
}

// TODO: Remove this file/package's dependency on app

// Create new HTTP request from the TestRequest properties.
func (tr *TestRequest) CreateRequest() *http.Request {
  // Create body from buffered data if data exists but body doesn't.
  if tr.Data != nil && tr.Body == nil {
    tr.Body, _ = tr.Data.AsBuffer()
  }

  // Create new HTTP request.
  req, err := http.NewRequest(tr.Method, app.Config.BaseRoute() + tr.Route, tr.Body)

  if err != nil {
    panic(fmt.Errorf("error creating new http request object: %s", err.Error()))
  }

  // Add Sweet Tea auth header if specified as an 'authed' request.
  if tr.Authed {
    req.Header.Set(app.Config.AuthHeaderName, app.Config.RestApiToken)
  }

  // Add Content-Type header based on which method is being used.
  var contentType string

  if tr.Method == "GET" || tr.Method == "DELETE" {
    contentType = "application/x-www-form-urlencoded"
  } else {
    contentType = "application/json"
  }

  req.Header.Set("Content-Type", contentType)

  // Set any extra headers provided.
  for k, v := range tr.ExtraHeaders {
    req.Header.Set(k, v)
  }

  return req
}