package routes

import (
  "testing"
  "github.com/stretchr/testify/assert"
)

func TestCreateUserHandler(t *testing.T) {
  res := TestRouter.JSONRequest("POST", UserRoute, nil, false)
  assert.Equal(t, 401, res.StatusCode())
}