package pl

import (
  "net/http"
  "encoding/json"
  "gopkg.in/validator.v2"
)

// ----------- POST /users -----------

type CreateUserPayload struct {
  CallingEmail    string `json:"calling_email" validate:"nonzero"`
  CallingPassword string `json:"calling_password" validate:"nonzero"`
  NewEmail        string `json:"new_email" validate:"nonzero"`
  NewPassword     string `json:"new_password" validate:"nonzero"`
  Admin           bool   `json:"admin"`
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

// ----------- POST /companies -----------

type CreateCompanyPayload struct {
  CallingEmail    string `json:"calling_email" validate:"nonzero"`
  CallingPassword string `json:"calling_password" validate:"nonzero"`
  Name            string `json:"name" validate:"nonzero"`
}

func (p *CreateCompanyPayload) Validate(req *http.Request) bool {
  return json.NewDecoder(req.Body).Decode(p) == nil && validator.Validate(p) == nil
}
