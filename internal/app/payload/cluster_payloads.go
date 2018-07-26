package payload

import (
  "net/http"
  "encoding/json"
  "gopkg.in/validator.v2"
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/pkg/util/cloud"
)

// ----------- POST /cluster -----------

type CreateClusterPayload struct {
  ExecutorEmail    string `json:"executor_email" validate:"nonzero"`
  ExecutorPassword string `json:"executor_password" validate:"nonzero"`
  Name             string `json:"name" validate:"nonzero"`
  Cloud            string `json:"cloud" validate:"nonzero"`
  State            string `json:"state"`
}

func (p *CreateClusterPayload) Validate(req *http.Request) bool {
  if err := json.NewDecoder(req.Body).Decode(p); err != nil {
    return false
  }

  if validator.Validate(p) != nil {
    return false
  }

  // State can't be empty unless on test or local environment.
  if p.State == "" && !app.Config.OnTest() && !app.Config.OnLocal() {
    return false
  }

  if !cloud.IsValidCloud(p.Cloud) {
    return false
  }

  return true
}

// ----------- PUT /cluster -----------

type UpdateClusterPayload struct {
  ExecutorEmail    string `json:"executor_email" validate:"nonzero"`
  ExecutorPassword string `json:"executor_password" validate:"nonzero"`
  Slug             string `json:"slug" validate:"nonzero"`
  Name             string `json:"name"`
  Cloud            string `json:"cloud"`
  State            string `json:"state"`
}

func (p *UpdateClusterPayload) Validate(req *http.Request) bool {
  return json.NewDecoder(req.Body).Decode(p) == nil && validator.Validate(p) == nil
}
