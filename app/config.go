package app

import (
  "fmt"
  "github.com/joeshaw/envdecode"
  "github.com/sweettea-io/rest-api/pkg/utils"
)

type appConfig struct {
  APIVersion          string  `env:"API_VERSION,required"`
  AWSAccessKeyId      string  `env:"AWS_ACCESS_KEY_ID,required"`
  AWSRegionName       string  `env:"AWS_REGION_NAME,required"`
  AWSSecretAccessKey  string  `env:"AWS_SECRET_ACCESS_KEY,required"`
  BuildClusterName    string  `env:"BUILD_CLUSTER_NAME,required"`
  BuildClusterState   string  `env:"BUILD_CLUSTER_STATE"`
  CloudProvider       string  `env:"CLOUD_PROVIDER,required"`
  DatabaseUrl         string  `env:"DATABASE_URL,required"`
  Debug               bool    `env:"DEBUG,required"`
  Domain              string  `env:"DOMAIN,required"`
  Env                 string  `env:"ENV,required"`
  HostedZoneId        string  `env:"HOSTED_ZONE_ID"`
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
  UserCreationHash    string  `env:"USER_CREATION_HASH"`
  WildcardSSLCertArn  string  `env:"WILDCARD_SSL_CERT_ARN"`
  WorkerCount         uint    `env:"WORKER_COUNT,required"`
}

var Config appConfig

func LoadConfig() {
  // Unmarshal envs into Config struct.
  if err := envdecode.Decode(&Config); err != nil {
    panic(fmt.Errorf("Failed to load app config: %s\n", err.Error()))
  }

  // --- Evaluate any configs reliant on more sophisticated validation ---

  // Non-local env checks.
  if Config.Env != utils.Envs.Local {
    errMsg := "Failed to load app config: %s required on non-local environments.\n"

    // Not using for-loop for the following in case a non-string env is needed here in the future.

    // BUILD_CLUSTER_STATE is required.
    if Config.BuildClusterState == "" {
      panic(fmt.Errorf(errMsg, "BUILD_CLUSTER_STATE"))
    }

    // HOSTED_ZONE_ID is required.
    if Config.HostedZoneId == "" {
      panic(fmt.Errorf(errMsg, "HOSTED_ZONE_ID"))
    }

    // WILDCARD_SSL_CERT_ARN is required.
    if Config.WildcardSSLCertArn == "" {
      panic(fmt.Errorf(errMsg, "WILDCARD_SSL_CERT_ARN"))
    }
  }

  // Ensure CLOUD_PROVIDER value is supported.
  if !utils.IsValidCloud(Config.CloudProvider) {
    panic(fmt.Errorf(
      "%s is not a valid cloud provider. Check 'pkg/utils/clouds.go' for a list of valid options.\n",
      Config.CloudProvider,
    ))
  }
}