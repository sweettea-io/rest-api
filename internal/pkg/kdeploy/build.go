package kdeploy

import (
  "fmt"
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/pkg/model"
  "github.com/sweettea-io/rest-api/internal/pkg/model/buildable"
  "github.com/sweettea-io/rest-api/internal/pkg/service/buildablesvc"
  "github.com/sweettea-io/rest-api/internal/pkg/service/projectsvc"
  "github.com/sweettea-io/rest-api/internal/pkg/util/cluster"
  "github.com/sweettea-io/rest-api/internal/pkg/util/image"
  "github.com/sweettea-io/rest-api/internal/pkg/util/maputil"
  "github.com/sweettea-io/rest-api/internal/pkg/util/timeutil"
  "github.com/sweettea-io/rest-api/internal/pkg/util/typeconvert"
  corev1 "k8s.io/api/core/v1"
  typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type Build struct {
  // Args
  ResourceID    uint
  TargetCluster string

  // Establish on Init
  Project       *model.Project
  Buildable     buildable.Buildable
  ClusterName   string
  DeployName    string
  Image         string
  ContainerName string

  // K8S resources
  Namespace     string
  Client        *typedcorev1.CoreV1Client
  VolumeMounts  []corev1.VolumeMount
  Volumes       []corev1.Volume
  Envs          []corev1.EnvVar
  Containers    []corev1.Container
  Pod           *corev1.Pod
}

func (b *Build) Init(args map[string]interface{}) error {
  // Decode/assign args to struct keys.
  b.ResourceID = args["resourceID"].(uint)
  b.TargetCluster = args["targetCluster"].(string)

  // Find Project by ID.
  project, err := projectsvc.FromID(args["projectID"].(uint))

  if err != nil {
    return err
  }

  // Find Buildable model from resource id & target cluster.
  resource, err := buildablesvc.FromID(b.ResourceID, b.TargetCluster)

  if err != nil {
    return err
  }

  b.Project = project
  b.Buildable = resource
  b.ClusterName = app.Config.BuildClusterName
  b.ContainerName = fmt.Sprintf("%s-%s-%s", b.TargetCluster, cluster.Build, project.Uid)
  b.DeployName = fmt.Sprintf("%s-%v", b.ContainerName, timeutil.MSSinceEpoch())
  b.Image = fmt.Sprintf("%s/%s", app.Config.DockerRegistryOrg, image.BuildServer)

  return nil
}

func (b *Build) Configure() error {
  // Configure k8s client.
  if err := b.makeClient(); err != nil {
    return err
  }

  // Build up k8s resources into the Pod we wish to deploy.
  b.makeVolumeMounts()
  b.makeVolumes()
  b.makeEnvs()
  b.makeContainers()
  b.makePod()

  return nil
}

func (b *Build) Perform() error {
  // Deploy the configured pod.
  if _, err := b.Client.Pods(b.Namespace).Create(b.Pod); err != nil {
    return fmt.Errorf("error performing build deploy: %s", err.Error())
  }

  return nil
}

func (b *Build) Watch() {

}

// FollowOnDeploy returns the KDeploy instance responsible for
// deploying to the target cluster of the Build deploy.
func (b *Build) FollowOnDeploy() KDeploy {
  switch b.TargetCluster {
  case cluster.Train:
    return &Train{}
  case cluster.Api:
    return &Api{}
  default:
    return nil
  }
}

func (b *Build) makeClient() error {
  // Configure CoreV1 client.
  client, nsp, err := ConfigureCoreV1(b.ClusterName)
  
  if err != nil {
    return err
  }
  
  // Store client and namespace on build object.
  b.Client = client
  b.Namespace = nsp
  
  return nil
}

func (b *Build) makeVolumeMounts() {
  b.VolumeMounts = VolumeMounts([]map[string]string{{
    "name": "dockersock",
    "path": "/var/run",
  }})
}

func (b *Build) makeVolumes() {
  b.Volumes = Volumes([]map[string]string{{
    "name": "dockersock",
    "path": "/var/run",
  }})
}

func (b *Build) makeEnvs() {
  envs := map[string]string{
    "BUILD_TARGET_ACCESS_TOKEN": b.Project.GetHost().GetToken(),
    "BUILD_TARGET_SHA": b.Buildable.GetCommit().Sha,
    "BUILD_TARGET_UID": b.Project.Uid,
    "BUILD_TARGET_URL": b.Project.Url(),
    "LOG_STREAM_KEY": b.Buildable.GetUid(),
    "DOCKER_REGISTRY_ORG": app.Config.DockerRegistryOrg,
    "DOCKER_REGISTRY_USERNAME": app.Config.DockerRegistryUsername,
    "DOCKER_REGISTRY_PASSWORD": app.Config.DockerRegistryPassword,
    "REDIS_POOL_MAX_ACTIVE": string(app.Config.RedisPoolMaxActive),
    "REDIS_POOL_MAX_IDLE": string(app.Config.RedisPoolMaxIdle),
    "REDIS_POOL_WAIT": typeconvert.BoolToStr(app.Config.RedisPoolWait),
    "REDIS_ADDRESS": app.Config.RedisAddress,
    "REDIS_PASSWORD": app.Config.RedisPassword,
    "TARGET_CLUSTER": b.TargetCluster,
  }

  // Add env info for all buildpacks.
  envs = maputil.MergeMaps(envs, app.Config.BuildpackEnvs())

  b.Envs = EnvVars(envs)
}

func (b *Build) makeContainers() {
  b.Containers = Containers([]map[string]interface{}{{
    "name":         b.ContainerName,
    "image":        b.Image,
    "envs":         b.Envs,
    "volumeMounts": b.VolumeMounts,
  }})
}

func (b *Build) makePod() {
  b.Pod = Pod(map[string]interface{}{
    "label": b.DeployName,
    "containers": b.Containers,
    "volumes": b.Volumes,
    "restart": corev1.RestartPolicyNever,
  })
}