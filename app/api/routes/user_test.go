package routes

import (
  "testing"
  "github.com/stretchr/testify/assert"
)

func TestCreateUserHandler(t *testing.T) {
  defer Teardown()
  res := router.JSONRequest("POST", UserRoute, nil, false)
  assert.Equal(t, 401, res.StatusCode())
}