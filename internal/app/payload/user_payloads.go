package payload

import (
  "net/http"
  "encoding/json"
  "gopkg.in/validator.v2"
)

// ----------- POST /user -----------

type CreateUserPayload struct {
  ExecutorEmail    string `json:"executor_email"`
  ExecutorPassword string `json:"executor_password" validate:"nonzero"`
  NewEmail         string `json:"new_email" validate:"nonzero"`
  NewPassword      string `json:"new_password" validate:"nonzero"`
  Admin            bool   `json:"admin"`
}

func (p *CreateUserPayload) Validate(req *http.Request) bool {
  return json.NewDecoder(req.Body).Decode(p) == nil && validator.Validate(p) == nil
}

// ----------- POST /user/auth -----------

type UserAuthPayload struct {
  Email    string `json:"email" validate:"nonzero"`
  Password string `json:"password" validate:"nonzero"`
}

func (p *UserAuthPayload) Validate(req *http.Request) bool {
  return json.NewDecoder(req.Body).Decode(p) == nil && validator.Validate(p) == nil
}