package kdeploy

import (
  "fmt"
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/pkg/model"
  "github.com/sweettea-io/rest-api/internal/pkg/service/trainjobsvc"
  "github.com/sweettea-io/rest-api/internal/pkg/util/cluster"
  "github.com/sweettea-io/rest-api/internal/pkg/util/maputil"
  "github.com/sweettea-io/rest-api/internal/pkg/util/timeutil"
  corev1 "k8s.io/api/core/v1"
  typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type Train struct {
  // Establish on Init
  CustomEnvs    map[string]string
  TrainJob      *model.TrainJob
  Commit        *model.Commit
  Project       *model.Project
  ModelVersion  *model.ModelVersion
  Model         *model.Model
  ClusterName   string
  DeployName    string
  Image         string
  ContainerName string

  // K8S resources
  Namespace     string
  Client        *typedcorev1.CoreV1Client
  Envs          []corev1.EnvVar
  Containers    []corev1.Container
  Pod           *corev1.Pod

}

func (t *Train) Init(args map[string]interface{}) error {
  // Set custom train deploy envs.
  customEnvs, ok := args["envs"].(map[string]string)
  if !ok {
    return fmt.Errorf("error parsing custom train deploy envs")
  }

  // Find TrainJob by ID.
  trainJob, err := trainjobsvc.FromID(args["trainJobID"].(uint))

  if err != nil {
    return err
  }

  t.CustomEnvs = customEnvs
  t.TrainJob = trainJob
  t.Commit = &trainJob.Commit
  t.Project = &trainJob.Commit.Project
  t.ModelVersion = &trainJob.ModelVersion
  t.Model = &trainJob.ModelVersion.Model

  // Cluster name is the SweetTea Train Cluster name.
  t.ClusterName = app.Config.TrainClusterName

  // Create unique container and deploy names.
  t.ContainerName = fmt.Sprintf("%s-%s", cluster.Train, t.Project.Uid)
  t.DeployName = fmt.Sprintf("%s-%v", t.ContainerName, timeutil.MSSinceEpoch())

  // Ex: sweetteaprod/train-<project_uid>:<commit_sha>
  t.Image = fmt.Sprintf("%s/%s:%s", app.Config.DockerRegistryOrg, t.ContainerName, t.Commit.Sha)

  return nil
}

func (t *Train) Configure() error {
  // Configure k8s client.
  if err := t.makeClient(); err != nil {
    return err
  }

  // Build up k8s resources into the Pod we wish to deploy.
  t.makeEnvs()
  t.makeContainers()
  t.makePod()

  return nil
}

// Perform deploys the configured pod to the Train Cluster.
func (t *Train) Perform() error {
  return DeployPod(t.Client, t.Namespace, t.Pod, cluster.Train)
}

func (t *Train) makeClient() error {
  // Configure CoreV1 client.
  client, nsp, err := ConfigureCoreV1(t.ClusterName)

  if err != nil {
    return err
  }

  // Store client and namespace on build object.
  t.Client = client
  t.Namespace = nsp

  return nil
}

func (t *Train) makeEnvs() {
  // Create envs from custom+deploy-agnostic envs.
  envs := maputil.MergeMaps(t.CustomEnvs, map[string]string{
    "AWS_ACCESS_KEY_ID": app.Config.AWSAccessKeyId,
    "AWS_REGION_NAME": app.Config.AWSRegionName,
    "AWS_SECRET_ACCESS_KEY": app.Config.AWSSecretAccessKey,
    "PROJECT_UID": t.Project.Uid,
    "REDIS_ADDRESS": app.Config.RedisAddress,
    "REDIS_PASSWORD": app.Config.RedisPassword,
  })

  t.Envs = EnvVars(envs)
}

func (t *Train) makeContainers() {
  t.Containers = Containers([]map[string]interface{}{{
    "name": t.ContainerName,
    "image": t.Image,
    "envs": t.Envs,
  }})
}

func (t *Train) makePod() {
  t.Pod = Pod(map[string]interface{}{
    "label": t.DeployName,
    "containers": t.Containers,
    "restart": corev1.RestartPolicyNever,
  })
}