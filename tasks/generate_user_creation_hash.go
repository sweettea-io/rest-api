package main

import (
  "fmt"
  "os"
  "github.com/sweettea-io/rest-api/internal/pkg/util/crypt"
)

func main() {
  // Ensure password arg was provided.
  if len(os.Args) <= 1 {
    panic(fmt.Errorf("task must be called with one argument -- the password to hash"))
  }

  // Hash the provided password.
  hash := crypt.Sha256Hash(os.Args[1])

  // Let the executor know what the hashed password is.
  fmt.Printf("Successfully generated hash.\nSet the USER_CREATION_HASH env var to %s for the environment you wish to use this for.\n", hash)
}