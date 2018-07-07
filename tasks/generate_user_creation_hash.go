package main

import (
  "github.com/sweettea-io/rest-api/pkg/utils"
  "fmt"
  "os"
)

func main() {
  if len(os.Args) <= 1 {
    panic(fmt.Errorf("Task must be called with one argument -- the password to hash."))
  }

  hash, err := utils.HashPw(os.Args[1])

  if err != nil {
    panic(fmt.Errorf("Error generating hash: %s", err.Error()))
  }

  fmt.Printf("Successfully generated hash.\nSet the USER_CREATION_HASH env var to %s for the environment you wish to use this for.\n", hash)
}