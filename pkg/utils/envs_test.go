package utils

import (
  "testing"
  "github.com/stretchr/testify/assert"
)

func TestIsValidEnv(t *testing.T) {
  testCases := []struct {
    name     string
    env      string
    expected bool
  }{
    {
      name: "invalid env",
      env: "invalid-env-value",
      expected: false,
    },
    {
      name: "test env",
      env: "test",
      expected: true,
    },
    {
      name: "local env",
      env: "local",
      expected: true,
    },
    {
      name: "dev env",
      env: "dev",
      expected: true,
    },
    {
      name: "staging env",
      env: "staging",
      expected: true,
    },
    {
      name: "prod env",
      env: "prod",
      expected: true,
    },
  }

  for _, tc := range testCases {
    result := IsValidEnv(tc.env)
    assert.Equal(t, tc.expected, result, tc.name)
  }
}