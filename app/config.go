package app

import (
  "github.com/kelseyhightower/envconfig"
  "fmt"
)

type appConfig struct {
  DatabaseUrl string `required:"true"`
  Debug       bool   `default:"false"`
  Env         string `default:"dev"`
  Port        int    `default:"5000"`
  ApiVersion  string `default:"v1"`
}

var Config appConfig

func LoadConfig() {
  // Look for env vars starting with "ST_".
  prefix := "st"

  // Unmarshal envs into Config struct.
  if err := envconfig.Process(prefix, &Config); err != nil {
    panic(fmt.Errorf("Failed to load app config: %s", err.Error()))
  }
}