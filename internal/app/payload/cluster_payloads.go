package payload

import (
  "net/http"
  "encoding/json"
  "gopkg.in/validator.v2"
)

// ----------- POST /clusters -----------

type CreateClusterPayload struct {
  ExecutorEmail    string `json:"executor_email" validate:"nonzero"`
  ExecutorPassword string `json:"executor_password" validate:"nonzero"`
  CompanyName      string `json:"company_name" validate:"nonzero"`
  Name             string `json:"name" validate:"nonzero"`
  Cloud            string `json:"cloud" validate:"nonzero"`
  State            string `json:"state"`
}

func (p *CreateClusterPayload) Validate(req *http.Request) bool {
  return json.NewDecoder(req.Body).Decode(p) == nil && validator.Validate(p) == nil
}