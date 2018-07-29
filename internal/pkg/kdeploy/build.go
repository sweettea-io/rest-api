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
  corev1 "k8s.io/api/core/v1"
  metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
  typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
  "github.com/sweettea-io/rest-api/internal/pkg/util/typeconvert"
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
  b.DeployName = fmt.Sprintf("%s-%v", b.ContainerName, 1)
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
  b.VolumeMounts = []corev1.VolumeMount{{
    Name: "dockersock",
    MountPath: "/var/run",
  }}
}

func (b *Build) makeVolumes() {
  b.Volumes = []corev1.Volume{
    CoreV1Volume("dockersock", "/var/run"),
  }
}

func (b *Build) makeEnvs() {
  envs := map[string]string{
    "BUILD_TARGET_ACCESS_TOKEN": "", // TODO: Figure out best way to address this.
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

  b.Envs = CoreV1EnvVars(envs)
}

func (b *Build) makeContainers() {
  b.Containers = []corev1.Container{{
    Name: b.ContainerName,
    Image: b.Image,
    Env: b.Envs,
    VolumeMounts: b.VolumeMounts,
  }}
}

func (b *Build) makePod() {
  b.Pod = &corev1.Pod{
    ObjectMeta: metav1.ObjectMeta{
      Labels: map[string]string{
        "app": b.DeployName,
      },
    },
    Spec: corev1.PodSpec{
      Containers: b.Containers,
      Volumes: b.Volumes,
      RestartPolicy: "Never",
    },
  }
}