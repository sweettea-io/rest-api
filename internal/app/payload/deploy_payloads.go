package payload

import (
  "net/http"
  "encoding/json"
  "gopkg.in/validator.v2"
  "github.com/sweettea-io/rest-api/internal/pkg/util/str"
  "github.com/sweettea-io/rest-api/internal/pkg/service/modelsvc"
)

// ----------- POST /deploy -----------

type CreateDeployPayload struct {
  Name           string `json:"name" validate:"nonzero"`
  ApiClusterName string `json:"apiCluster" validate:"nonzero"`
  ProjectNsp     string `json:"projectNsp" validate:"nonzero"`
  Model          string `json:"model"`
  Sha            string `json:"sha"`
  Envs           string `json:"envs"`
}

func (p *CreateDeployPayload) ModelBreakdown() (string, string) {
  return modelsvc.Breakdown(p.Model)
}

func (p *CreateDeployPayload) ApiClusterSlug() string {
  return str.Slugify(p.ApiClusterName)
}

func (p *CreateDeployPayload) Validate(req *http.Request) bool {
  return json.NewDecoder(req.Body).Decode(p) == nil && validator.Validate(p) == nil
}

// ----------- PUT /deploy -----------

type UpdateDeployPayload struct {
  Name       string `json:"name" validate:"nonzero"`
  ProjectNsp string `json:"projectNsp" validate:"nonzero"`
  Model      string `json:"model"`
  Sha        string `json:"sha"`
  Envs       string `json:"envs"`
}

func (p *UpdateDeployPayload) ModelBreakdown() (string, string) {
  return modelsvc.Breakdown(p.Model)
}

func (p *UpdateDeployPayload) Slug() string {
  return str.Slugify(p.Name)
}

func (p *UpdateDeployPayload) Validate(req *http.Request) bool {
  return json.NewDecoder(req.Body).Decode(p) == nil && validator.Validate(p) == nil
}