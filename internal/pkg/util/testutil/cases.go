package testutil

import "github.com/sweettea-io/rest-api/internal/pkg/util/enc"

type RequestCase struct {
  Name             string
  Request          *Request
  ExpectedStatus   int
  ExpectedRespJSON *enc.JSON
}