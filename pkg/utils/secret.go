package utils

import (
  "encoding/base64"
  "crypto/rand"
)

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomBytes(n int) ([]byte, error) {
  b := make([]byte, n)
  _, err := rand.Read(b)
  // Note that err == nil only if we read len(b) bytes.
  if err != nil {
    return nil, err
  }

  return b, nil
}

// GenerateRandomStringURLSafe returns a URL-safe, base64 encoded securely generated random string.
func GenerateRandomStringURLSafe(n int) (string, error) {
  b, err := GenerateRandomBytes(n)
  return base64.URLEncoding.EncodeToString(b), err
}

// Mint a new secret token.
func FreshSecret() (string, error) {
  return GenerateRandomStringURLSafe(32)
}