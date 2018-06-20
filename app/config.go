package app

import (
  "fmt"
)

type appConfig struct {
  // Required (must exist as env var)
  AwsAccessKeyId     string `env:"AWS_ACCESS_KEY_ID,required"`
  AwsRegionName      string `env:"AWS_REGION_NAME,required"`
  AwsSecretAccessKey string `env:"AWS_SECRET_ACCESS_KEY,required"`
  BuildClusterName   string `env:"BUILD_CLUSTER_NAME,required"`
  BuildClusterState  string `env:"BUILD_CLUSTER_STATE,required"`
  CoreClusterName    string `env:"CORE_CLUSTER_NAME,required"`
  CoreClusterState   string `env:"CORE_CLUSTER_STATE,required"`
  DatabaseUrl        string `env:"DATABASE_URL,required"`
  Domain             string `env:"DOMAIN,required"`
  Env                string `env:"ENV,required"`
  RedisUrl           string `env:"REDIS_URL,required"`

  // Optional with defaults
  ApiClusterZones    string `env:"API_CLUSTER_ZONES,default=us-west-1a"`
  ApiVersion         string `env:"API_VERSION,default=v1"`
  ClusterImage       string `env:"CLUSTER_IMAGE,default=099720109477/ubuntu/images/hvm-ssd/ubuntu-xenial-16.04-amd64-server-20171026.1"`
  Debug              bool   `env:"DEBUG,default=true"`

  JobQueueNsp        string `default:"st_job_queue"`
  ServerPort               int    `default:"80"`
  RedisPoolMaxActive int    `default:"5"`
  RedisPoolMaxIdle   int    `default:"5"`
  RedisPoolWait      bool   `default:"true"`
  WorkerCount        uint   `default:"10"`
}

var Config appConfig

func LoadConfig() {
  // Look for env vars starting with "ST_".
  prefix := "st"

  // Unmarshal envs into Config struct.
  //if err := envconfig.Process(prefix, &Config); err != nil {
  //  panic(fmt.Errorf("Failed to load app config: %s", err.Error()))
  //}
}