package config

import (
  "fmt"
  "github.com/sweettea-io/envdecode"
  "github.com/sweettea-io/rest-api/internal/pkg/util/cloud"
  "github.com/sweettea-io/rest-api/internal/pkg/util/env"
)

// Config represents app config populated from environment variables.
type Config struct {
  APIVersion             string `env:"API_VERSION,required"`
  AuthHeaderName         string `env:"AUTH_HEADER_NAME,required"`
  AWSAccessKeyId         string `env:"AWS_ACCESS_KEY_ID,required,ignore_on_envs=test|local"`
  AWSRegionName          string `env:"AWS_REGION_NAME,required,ignore_on_envs=test|local"`
  AWSSecretAccessKey     string `env:"AWS_SECRET_ACCESS_KEY,required,ignore_on_envs=test|local"`
  BuildClusterName       string `env:"BUILD_CLUSTER_NAME,required,ignore_on_envs=test"`
  BuildClusterState      string `env:"BUILD_CLUSTER_STATE,required,ignore_on_envs=test|local"`
  CloudProvider          string `env:"CLOUD_PROVIDER,required,ignore_on_envs=test|local"`
  DatabaseUrl            string `env:"DATABASE_URL,required"`
  Debug                  bool   `env:"DEBUG,required"`
  DockerRegistryOrg      string `env:"DOCKER_REGISTRY_ORG,required,ignore_on_envs=test"`
  DockerRegistryUsername string `env:"DOCKER_REGISTRY_USERNAME,required,ignore_on_envs=test"`
  DockerRegistryPassword string `env:"DOCKER_REGISTRY_PASSWORD,required,ignore_on_envs=test"`
  Domain                 string `env:"DOMAIN,required,ignore_on_envs=test"`
  Env                    string `env:"ENV,required"`
  HostedZoneId           string `env:"HOSTED_ZONE_ID,ignore_on_envs=test|local"`
  JobQueueNsp            string `env:"JOB_QUEUE_NSP,required"`
  KubeConfig             string `env:"KUBECONFIG,required"`
  RedisPoolMaxActive     int    `env:"REDIS_POOL_MAX_ACTIVE,required"`
  RedisPoolMaxIdle       int    `env:"REDIS_POOL_MAX_IDLE,required"`
  RedisPoolWait          bool   `env:"REDIS_POOL_WAIT,required"`
  RedisAddress           string `env:"REDIS_ADDRESS,required"`
  RedisPassword          string `env:"REDIS_PASSWORD"`
  RestApiToken           string `env:"REST_API_TOKEN,required"`
  ServerPort             int    `env:"SERVER_PORT,required"`
  TrainClusterName       string `env:"TRAIN_CLUSTER_NAME"`
  TrainClusterState      string `env:"TRAIN_CLUSTER_STATE"`
  UserCreationHash       string `env:"USER_CREATION_HASH"`
  WildcardSSLCertArn     string `env:"WILDCARD_SSL_CERT_ARN,required,ignore_on_envs=test|local"`
  WorkerCount            uint   `env:"WORKER_COUNT,required,ignore_on_envs=test"`
}

// BaseRoute returns the base route for the server app's API.
// The base route is created from the API version config value.
func (cfg *Config) BaseRoute() string {
  return fmt.Sprintf("/%s", cfg.APIVersion)
}

// New creates and returns a new Config struct instance populated from environment variables.
func New() *Config {
  // Unmarshal values into a config struct.
  var cfg Config
  if err := envdecode.Decode(&cfg); err != nil {
    panic(fmt.Errorf("Failed to load app config: %s\n", err.Error()))
  }

  // Validate config values.
  validateConfigs(&cfg)

  return &cfg
}

// Validate Config values even further than what has
// already been done during the Decode process.
func validateConfigs(cfg *Config) {
  // Ensure ENV value is supported.
  if !env.IsValidEnv(cfg.Env) {
    panic(fmt.Errorf(
      "%s is not a valid env. Check 'internal/pkg/util/env/tiers.go' for a list of valid options.\n",
      cfg.Env,
    ))
  }

  // Ensure CLOUD_PROVIDER value is supported (if it exists -- not always required)
  if cfg.CloudProvider != "" && !cloud.IsValidCloud(cfg.CloudProvider) {
    panic(fmt.Errorf(
      "%s is not a valid cloud provider. Check 'internal/pkg/util/cloud/clouds.go' for a list of valid options.\n",
      cfg.CloudProvider,
    ))
  }
}