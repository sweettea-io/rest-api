package payload

import (
  "net/http"
  "encoding/json"
  "gopkg.in/validator.v2"
  "github.com/sweettea-io/rest-api/internal/pkg/util/str"
)

// ----------- POST /train_job -----------

type CreateTrainJobPayload struct {
  ProjectNsp string `json:"projectNsp" validate:"nonzero"`
  ModelName  string `json:"modelName" default:"default"`
  Envs       string `json:"envs"`
}

func (p *CreateTrainJobPayload) ModelSlug() string {
  return str.Slugify(p.ModelName)
}

func (p *CreateTrainJobPayload) Validate(req *http.Request) bool {
  return json.NewDecoder(req.Body).Decode(p) == nil && validator.Validate(p) == nil
}
