package utils

import "github.com/satori/go.uuid"

func NewUid() string {
  return uuid.Must(uuid.NewV4(), nil).String()
}