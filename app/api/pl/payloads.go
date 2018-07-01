package pl

import (
  "net/http"
  "encoding/json"
  "gopkg.in/validator.v2"
)

// ----------- POST /users -----------

type CreateUserPayload struct {
  Email    string `json:"email" validate:"nonzero"`
  Password string `json:"password" validate:"nonzero"`
  Admin    bool   `json:"admin"`
}

func (p *CreateUserPayload) Validate(req *http.Request) bool {
  return json.NewDecoder(req.Body).Decode(p) == nil && validator.Validate(p) == nil
}

// ----------- POST /users/auth -----------

type UserAuthPayload struct {
  Email    string `json:"email" validate:"nonzero"`
  Password string `json:"password" validate:"nonzero"`
}

func (p *UserAuthPayload) Validate(req *http.Request) bool {
  return json.NewDecoder(req.Body).Decode(p) == nil && validator.Validate(p) == nil
}