package cases

import (
  "github.com/sweettea-io/rest-api/pkg/utils"
  "github.com/sweettea-io/rest-api/pkg/test_utils"
)

type RequestCase struct {
  Name             string
  Request          *test_utils.TestRequest
  ExpectedStatus   int
  ExpectedRespJSON *utils.JSON
}