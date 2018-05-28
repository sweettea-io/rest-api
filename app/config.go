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
  if err := envconfig.Process("st", &Config); err != nil {
    return err
  }

  return nil
}