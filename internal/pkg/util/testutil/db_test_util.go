package testutil

import (
  "testing"
  "github.com/stretchr/testify/assert"
  "github.com/sweettea-io/rest-api/internal/app"
)

func AssertTableCount(t *testing.T, tc *RequestCase, tableName string, count int) {
  // TODO: probably check for SQL injection here with tableName.
  // Get count of records in table by provided name.
  var actualCount int
  app.DB.Table(tableName).Count(&actualCount)

  // Assert counts are the same.
  assert.Equal(t, count, actualCount, tc.Name)
}
