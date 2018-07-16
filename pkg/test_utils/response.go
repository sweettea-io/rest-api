package test_utils

import (
  "net/http/httptest"
  "encoding/json"
  "github.com/sweettea-io/rest-api/pkg/utils"
)

// Wrapper type around httptest.ResponseRecorder, providing additional helper functions for
// JSON parsing the response body and fetching the response status code.
type TestResponse struct {
  raw *httptest.ResponseRecorder
}

// Parse and return the response body as utils.JSON type.
func (res *TestResponse) ParseJSON() *utils.JSON {
  // Parse body into JSON type.
  var data utils.JSON
  json.Unmarshal(res.raw.Body.Bytes(), &data)

  return &data
}

// Parse and return the response status code.
func (res *TestResponse) StatusCode() int {
  return res.raw.Result().StatusCode
}

