package payload

import (
  "net/http"
  "encoding/json"
  "gopkg.in/validator.v2"
  "strings"
  "github.com/sweettea-io/rest-api/internal/pkg/util/str"
  "github.com/sweettea-io/rest-api/internal/pkg/model"
)

// ----------- POST /deploy -----------

type CreateDeployPayload struct {
  ApiClusterName string `json:"apiCluster" validate:"nonzero"`
  ProjectNsp     string `json:"projectNsp" validate:"nonzero"`
  Model          string `json:"model"`
  Sha            string `json:"sha"`
  Envs           string `json:"envs"`
}

func (p *CreateDeployPayload) ModelBreakdown() (string, string) {
  if p.Model == "" {
    return model.DefaultModelSlug, ""
  }

  breakdown := strings.Split(p.Model, ":")
  slug := str.Slugify(breakdown[0])
  version := ""

  if len(breakdown) > 1 {
    version = breakdown[1]
  }

  return slug, version
}

func (p *CreateDeployPayload) ApiClusterSlug() string {
  return str.Slugify(p.ApiClusterName)
}

func (p *CreateDeployPayload) Validate(req *http.Request) bool {
  return json.NewDecoder(req.Body).Decode(p) == nil && validator.Validate(p) == nil
}
