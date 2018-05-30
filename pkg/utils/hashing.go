package utils

import "golang.org/x/crypto/bcrypt"

// Hash a password string using bcrypt.
func HashPw(password string) (string, error) {
  bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
  return string(bytes), err
}

// Verify a password hashes to the provided hash using bcrypt.
func VerifyPw(password, hash string) bool {
  err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
  return err == nil
}