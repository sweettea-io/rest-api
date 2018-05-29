package app

import (
  "github.com/kelseyhightower/envconfig"
)

type appConfig struct {
  DatabaseUrl string `required:"true"`
  Debug       bool   `default:"false"`
  Env         string `default:"dev"`
  Port        int    `default:"5000"`
  ApiVersion  string `default:"v1"`
}

var Config appConfig

func LoadConfig() error {
  // Look for env vars starting with "ST_".
  envVarPrefix := "st"

  // Unmarshal envs into Config struct.
  if err := envconfig.Process(envVarPrefix, &Config); err != nil {
    return err
  }

  return nil
}