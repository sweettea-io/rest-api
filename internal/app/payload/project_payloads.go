package payload

import (
  "net/http"
  "encoding/json"
  "gopkg.in/validator.v2"
  "github.com/sweettea-io/rest-api/internal/pkg/service/projectsvc"
)

// ----------- POST /project -----------

type CreateProjectPayload struct {
  Nsp string `json:"nsp" validate:"nonzero"`
}

func (p *CreateProjectPayload) Validate(req *http.Request) bool {
  return req.Body != nil &&
    json.NewDecoder(req.Body).Decode(p) == nil &&
    validator.Validate(p) == nil &&
    projectsvc.IsValidNsp(p.Nsp)
}

// ----------- GET /project -----------

type GetProjectPayload struct {
  Nsp string `json:"nsp"`
}

func (p *GetProjectPayload) Validate(req *http.Request) bool {
  return req.Body != nil && json.NewDecoder(req.Body).Decode(p) == nil && validator.Validate(p) == nil
}

// ----------- DELETE /project -----------

type DeleteProjectPayload struct {
  Nsp string `json:"nsp" validate:"nonzero"`
}

func (p *DeleteProjectPayload) Validate(req *http.Request) bool {
  return req.Body != nil && json.NewDecoder(req.Body).Decode(p) == nil && validator.Validate(p) == nil
}