package utils

import (
  "testing"
  "github.com/stretchr/testify/assert"
)

func TestVerifyBcrypt(t *testing.T) {
  blueHash, _ := BcryptHash("blue")

  // Test passwords do NOT match.
  assert.Equal(t, false, VerifyBcrypt("red", blueHash), "passwords do NOT match")

  // Test passwords do match.
  assert.Equal(t, true, VerifyBcrypt("blue", blueHash), "passwords do match")
}

func TestVerifySha256(t *testing.T) {
  blueHash := Sha256Hash("blue")

  // Test passwords do NOT match.
  assert.Equal(t, false, VerifySha256("red", blueHash), "passwords do NOT match")

  // Test passwords do match.
  assert.Equal(t, true, VerifySha256("blue", blueHash), "passwords do match")
}