package app

import (
  "github.com/kelseyhightower/envconfig"
  "fmt"
)

type appConfig struct {
  // Required (check env vars)
  DatabaseUrl        string `required:"true"`
  RedisUrl           string `required:"true"`

  // Optional with defaults
  ApiVersion         string `default:"v1"`
  Debug              bool   `default:"false"`
  Env                string `default:"dev"`
  JobQueueNsp        string `default:"st_job_queue"`
  Port               int    `default:"5000"`
  RedisPoolMaxActive int    `default:"5"`
  RedisPoolMaxIdle   int    `default:"5"`
  RedisPoolWait      bool   `default:"true"`
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