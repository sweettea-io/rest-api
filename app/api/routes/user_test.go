package routes

import (
  "testing"
  "github.com/stretchr/testify/assert"
)

func TestCreateUserHandler(t *testing.T) {
  // TODO: Add this to TestMain(m *testing.M)
  ClearTables()

  res := TestRouter.JSONRequest("POST", UserRoute, nil, false)

  assert.Equal(t, 200, res.StatusCode())
}