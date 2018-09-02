package testutil

import (
  "github.com/sweettea-io/rest-api/internal/pkg/util/enc"
  "github.com/sweettea-io/rest-api/internal/pkg/config"
  "testing"
)

type RequestCase struct {
  Name             string
  Request          *Request
  ExpectedStatus   int
  ExpectedRespJSON *enc.JSON
  CustomCfg        config.ConfigItf
  CustomAssertions []CustomReqAssertion
}

type CustomReqAssertion func(*testing.T, *RequestCase, int, *enc.JSON)

func (rc *RequestCase) SetupArgs() config.ConfigItf {
  return rc.CustomCfg
}