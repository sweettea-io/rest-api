package app

import (
  "fmt"
  "github.com/joeshaw/envdecode"
)

type appConfig struct {
  ApiClusterZones     string  `env:"API_CLUSTER_ZONES"`
  ApiVersion          string  `env:"API_VERSION,required"`
  AwsAccessKeyId      string  `env:"AWS_ACCESS_KEY_ID,required"`
  AwsRegionName       string  `env:"AWS_REGION_NAME,required"`
  AwsSecretAccessKey  string  `env:"AWS_SECRET_ACCESS_KEY,required"`
  BuildClusterName    string  `env:"BUILD_CLUSTER_NAME,required"`
  BuildClusterState   string  `env:"BUILD_CLUSTER_STATE,required"`
  CloudProvider       string  `env:"CLOUD_PROVIDER,required"`
  ClusterImage        string  `env:"CLUSTER_IMAGE,required"`
  CoreClusterName     string  `env:"CORE_CLUSTER_NAME,required"`
  CoreClusterState    string  `env:"CORE_CLUSTER_STATE,required"`
  DatabaseUrl         string  `env:"DATABASE_URL,required"`
  Debug               bool    `env:"DEBUG,required"`
  Domain              string  `env:"DOMAIN,required"`
  Env                 string  `env:"ENV,required"`
  HostedZoneId        string  `env:"HOSTED_ZONE_ID,required"`
  ImageOwner          string  `env:"IMAGE_OWNER,required"`
  ImageOwnerPw        string  `env:"IMAGE_OWNER_PW,required"`
  JobQueueNsp         string  `env:"JOB_QUEUE_NSP,required"`
  K8sVersion          string  `env:"K8S_VERSION,required"`
  KubeConfig          string  `env:"KUBECONFIG,required"`
  MasterSize          string  `env:"MASTER_SIZE,required"`
  NodeCount           int     `env:"NODE_COUNT,required"`
  NodeSize            string  `env:"NODE_SIZE,required"`
  RedisPoolMaxActive  int     `env:"REDIS_POOL_MAX_ACTIVE,required"`
  RedisPoolMaxIdle    int     `env:"REDIS_POOL_MAX_IDLE,required"`
  RedisPoolWait       bool    `env:"REDIS_POOL_WAIT,required"`
  RedisUrl            string  `env:"REDIS_URL,required"`
  RestApiToken        string  `env:"REST_API_TOKEN,required"`
  ServerPort          int     `env:"SERVER_PORT,required"`
  TrainClusterName    string  `env:"TRAIN_CLUSTER_NAME"`
  TrainClusterState   string  `env:"TRAIN_CLUSTER_STATE"`
  WildcardSslCertArn  string  `env:"WILDCARD_SSL_CERT_ARN,required"`
  WorkerCount         uint    `env:"WORKER_COUNT"`
}

var Config appConfig

func LoadConfig() {
  // Unmarshal envs into Config struct.
  if err := envdecode.Decode(&Config); err != nil {
   panic(fmt.Errorf("Failed to load app config: %s", err.Error()))
  }
}