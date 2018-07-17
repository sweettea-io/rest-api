package payload

import (
  "net/http"
  "encoding/json"
  "gopkg.in/validator.v2"
)

// ----------- POST /companies -----------

type CreateCompanyPayload struct {
  ExecutorEmail    string `json:"executor_email" validate:"nonzero"`
  ExecutorPassword string `json:"executor_password" validate:"nonzero"`
  Name             string `json:"name" validate:"nonzero"`
}

func (p *CreateCompanyPayload) Validate(req *http.Request) bool {
  return json.NewDecoder(req.Body).Decode(p) == nil && validator.Validate(p) == nil
}