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
  Slug             string `json:"name" validate:"nonzero"`
  Updates struct {
    Name  *string `json:"name"`
    Cloud *string `json:"cloud"`
    State *string `json:"state"`
  } `json:"updates"`
}

func (p *UpdateClusterPayload) Validate(req *http.Request) bool {
  if err := json.NewDecoder(req.Body).Decode(p); err != nil {
    return false
  }

  if validator.Validate(p) != nil {
    return false
  }

  // Validate state value.
  if p.Updates.State != nil && // if state was provided...
    *p.Updates.State == "" &&  // and it was provided as an empty string...
    !app.Config.OnTest() && !app.Config.OnLocal() { // the env must either be test or local to proceed...
    return false
  }

  // Validate cloud value.
  if !cloud.IsValidCloud(*p.Updates.Cloud) {
    return false
  }

  return true
}

func (p *UpdateClusterPayload) GetUpdates() *map[string]interface{} {
  updates := make(map[string]interface{})

  if p.Updates.Name != nil {
    updates["name"] = *p.Updates.Name
  }

  if p.Updates.Cloud != nil {
    updates["cloud"] = *p.Updates.Cloud
  }

  if p.Updates.State != nil {
    updates["state"] = *p.Updates.State
  }

  return &updates
}