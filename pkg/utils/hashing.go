package utils

import (
  "crypto/sha256"
  "encoding/base64"
  "golang.org/x/crypto/bcrypt"
)

// Hash a password string using bcrypt.
func BcryptHash(password string) (string, error) {
  bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
  return string(bytes), err
}

// Verify a password hashes to the provided hash using bcrypt.
func VerifyBcrypt(password, hash string) bool {
  err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
  return err == nil
}

// Hash a password string using sha256.
func Sha256Hash(password string) string {
  h := sha256.Sum256([]byte(password))
  return "{SHA256}" + base64.StdEncoding.EncodeToString(h[:])
}

// Verify a password hashes to the provided hash using sha256.
func VerifySha256(password, hash string) bool {
  return Sha256Hash(password) == hash
}