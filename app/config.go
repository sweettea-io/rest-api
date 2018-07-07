package app

import (
  "fmt"
  "github.com/joeshaw/envdecode"
  "github.com/sweettea-io/rest-api/pkg/utils"
)

type appConfig struct {
  ApiVersion          string  `env:"API_VERSION,required"`
  AwsAccessKeyId      string  `env:"AWS_ACCESS_KEY_ID,required"`
  AwsRegionName       string  `env:"AWS_REGION_NAME,required"`
  AwsSecretAccessKey  string  `env:"AWS_SECRET_ACCESS_KEY,required"`
  BuildClusterName    string  `env:"BUILD_CLUSTER_NAME,required"`
  BuildClusterState   string  `env:"BUILD_CLUSTER_STATE"`
  CloudProvider       string  `env:"CLOUD_PROVIDER,required"`
  DatabaseUrl         string  `env:"DATABASE_URL,required"`
  Debug               bool    `env:"DEBUG,required"`
  Domain              string  `env:"DOMAIN,required"`
  Env                 string  `env:"ENV,required"`
  HostedZoneId        string  `env:"HOSTED_ZONE_ID,required"`
  ImageOwner          string  `env:"IMAGE_OWNER,required"`
  ImageOwnerPw        string  `env:"IMAGE_OWNER_PW,required"`
  JobQueueNsp         string  `env:"JOB_QUEUE_NSP,required"`
  KubeConfig          string  `env:"KUBECONFIG,required"`
  RedisPoolMaxActive  int     `env:"REDIS_POOL_MAX_ACTIVE,required"`
  RedisPoolMaxIdle    int     `env:"REDIS_POOL_MAX_IDLE,required"`
  RedisPoolWait       bool    `env:"REDIS_POOL_WAIT,required"`
  RedisUrl            string  `env:"REDIS_URL,required"`
  RestApiToken        string  `env:"REST_API_TOKEN,required"`
  ServerPort          int     `env:"SERVER_PORT,required"`
  TrainClusterName    string  `env:"TRAIN_CLUSTER_NAME"`
  TrainClusterState   string  `env:"TRAIN_CLUSTER_STATE"`
  UserCreationHash    string  `env:"USER_CREATION_HASH,required"`
  WildcardSSLCertArn  string  `env:"WILDCARD_SSL_CERT_ARN,required"`
  WorkerCount         uint    `env:"WORKER_COUNT,required"`
}

var Config appConfig

func LoadConfig() {
  // Unmarshal envs into Config struct.
  if err := envdecode.Decode(&Config); err != nil {
    panic(fmt.Errorf("Failed to load app config: %s\n", err.Error()))
  }

  // --- Evaluate any configs reliant on more sophisticated validation ---

  // Ensure BuildClusterState exists for all non-local environments.
  if Config.BuildClusterState == "" && Config.Env != utils.Envs.Local {
    panic(fmt.Errorf("Failed to load app config: BUILD_CLUSTER_STATE required on non-local environments.\n"))
  }
}