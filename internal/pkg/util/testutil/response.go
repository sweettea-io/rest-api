package testutil

import (
  "encoding/json"
  "net/http/httptest"
  "github.com/sweettea-io/rest-api/internal/pkg/util/enc"
)

// Wrapper type around httptest.ResponseRecorder, providing additional helper functions for
// JSON parsing the response body and fetching the response status code.
type Response struct {
  raw *httptest.ResponseRecorder
}

// Parse and return the response body as utils.JSON type.
func (res *Response) ParseJSON() *enc.JSON {
  // Parse body into JSON type.
  var data enc.JSON
  json.Unmarshal(res.raw.Body.Bytes(), &data)

  return &data
}

// Parse and return the response status code.
func (res *Response) StatusCode() int {
  return res.raw.Result().StatusCode
}

