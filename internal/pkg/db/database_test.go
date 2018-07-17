package db

import "testing"

func TestConnection(t *testing.T) {
  defer func() {
    if r:= recover(); r == nil {
      t.Errorf("connection did not panic for bad DB url")
    }
  }()

  // Connection should fail if it can't connect to the database url.
  NewConnection("bad-db-url", true)
}