package payload

import (
  "net/http"
  "encoding/json"
  "gopkg.in/validator.v2"
)

// ----------- POST /projects -----------

type CreateProjectPayload struct {
  Name    string `json:"name" validate:"nonzero"`
  Company string `json:"company" validate:"nonzero"` // Company name (will be slugified and queried by this slug)
}

func (p *CreateProjectPayload) Validate(req *http.Request) bool {
  return json.NewDecoder(req.Body).Decode(p) == nil && validator.Validate(p) == nil
}