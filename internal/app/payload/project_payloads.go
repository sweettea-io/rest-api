package payload

import (
  "net/http"
  "encoding/json"
  "gopkg.in/validator.v2"
)

// ----------- POST /project -----------

type CreateProjectPayload struct {
  Name    string `json:"name" validate:"nonzero"`
}

func (p *CreateProjectPayload) Validate(req *http.Request) bool {
  return json.NewDecoder(req.Body).Decode(p) == nil && validator.Validate(p) == nil
}