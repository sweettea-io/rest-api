package pl

import (
  "net/http"
  "encoding/json"
  "gopkg.in/validator.v2"
)

// TODO: Add regex for emails in tags

// ----------- POST /users -----------

type CreateUserPayload struct {
  ExecutorEmail    string `json:"executor_email" validate:"nonzero"`
  ExecutorPassword string `json:"executor_password" validate:"nonzero"`
  NewEmail         string `json:"new_email" validate:"nonzero"`
  NewPassword      string `json:"new_password" validate:"nonzero"`
  Admin            bool   `json:"admin"`
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
  ExecutorEmail    string `json:"executor_email" validate:"nonzero"`
  ExecutorPassword string `json:"executor_password" validate:"nonzero"`
  Name             string `json:"name" validate:"nonzero"`
}

func (p *CreateCompanyPayload) Validate(req *http.Request) bool {
  return json.NewDecoder(req.Body).Decode(p) == nil && validator.Validate(p) == nil
}

// ----------- POST /clusters -----------

type CreateClusterPayload struct {
  ExecutorEmail    string `json:"executor_email" validate:"nonzero"`
  ExecutorPassword string `json:"executor_password" validate:"nonzero"`
  CompanyName      string `json:"company_name" validate:"nonzero"`
  Name             string `json:"name" validate:"nonzero"`
  Cloud            string `json:"cloud" validate:"nonzero"`
  State            string `json:"state" validate:"nonzero"`
}

func (p *CreateClusterPayload) Validate(req *http.Request) bool {
  return json.NewDecoder(req.Body).Decode(p) == nil && validator.Validate(p) == nil
}