package utils

import (
  "testing"
  "github.com/stretchr/testify/assert"
)

func TestIsValidCloud(t *testing.T) {
  testCases := []struct {
    name string
    cloud string
    expected bool
  }{
    {
      name: "invalid cloud",
      cloud: "invalid-cloud-value",
      expected: false,
    },
    {
      name: "valid cloud",
      cloud: "aws",
      expected: true,
    },
  }

  for _, tc := range testCases {
    result := IsValidCloud(tc.cloud)
    assert.Equal(t, tc.expected, result, tc.name)
  }
}