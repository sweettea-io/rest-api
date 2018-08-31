package testutil

import (
  "io"
  "net/http"
  "fmt"
  "github.com/sweettea-io/rest-api/internal/pkg/util/enc"
)

type Request struct {
  Method     string
  Route      string
  Body       io.Reader
  Data       *enc.JSON
  Headers    map[string]string
  BeforeSend []RequestModifier
}

type RequestModifier func(req *Request) (*Request, error)

// Create HTTP request out of the Request object.
func (req *Request) CreateHTTPRequest(baseRoute string) *http.Request {
  // Create body from buffered data if data exists but body doesn't.
  if req.Data != nil && req.Body == nil {
    req.Body, _ = req.Data.AsBuffer()
  }

  // Create new HTTP request.
  httpReq, err := http.NewRequest(req.Method, baseRoute + req.Route, req.Body)
  if err != nil {
    panic(fmt.Errorf("error creating new http request object: %s", err.Error()))
  }

  // Add Content-Type header based on which method is being used.
  var contentType string

  if req.Method == "GET" || req.Method == "DELETE" {
    contentType = "application/x-www-form-urlencoded"
  } else {
    contentType = "application/json"
  }

  httpReq.Header.Set("Content-Type", contentType)

  // Set any extra headers provided.
  for k, v := range req.Headers {
    httpReq.Header.Set(k, v)
  }

  return httpReq
}

func (req *Request) SetHeader(name string, value string) {
  if req.Headers == nil {
    req.Headers = make(map[string]string)
  }

  req.Headers[name] = value
}