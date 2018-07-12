package database

import "testing"

// Connection should fail if it can't connect to the database url.
func TestConnectionPanicsIfCantConnect(t *testing.T) {
  defer func() {
    if r:= recover(); r == nil {
      t.Errorf("Connection did not panic for bad DB url.")
    }
  }()

  Connection("bad-db-url")
}