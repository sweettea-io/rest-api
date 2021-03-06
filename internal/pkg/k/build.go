package k

import (
  "fmt"
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/pkg/model"
  "github.com/sweettea-io/rest-api/internal/pkg/service/buildablesvc"
  "github.com/sweettea-io/rest-api/internal/pkg/service/projectsvc"
  "github.com/sweettea-io/rest-api/internal/pkg/util/cluster"
  "github.com/sweettea-io/rest-api/internal/pkg/util/image"
  "github.com/sweettea-io/rest-api/internal/pkg/util/maputil"
  "github.com/sweettea-io/rest-api/internal/pkg/util/timeutil"
  "github.com/sweettea-io/rest-api/internal/pkg/util/typeconvert"
  "k8s.io/apimachinery/pkg/watch"
  corev1 "k8s.io/api/core/v1"
  typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type Build struct {
  // Args
  ResourceID       uint
  BuildTargetSha   string
  TargetCluster    string
  LogStreamKey     string

  // Establish on Init
  Project          *model.Project
  Buildable        model.Buildable
  ClusterName      string
  DeployName       string
  Image            string
  ContainerName    string
  ResultChannel    chan Result

  // K8S resourcess
  Namespace        string
  Client           *typedcorev1.CoreV1Client
  VolumeMounts     []corev1.VolumeMount
  Volumes          []corev1.Volume
  Envs             []corev1.EnvVar
  Containers       []corev1.Container
  Pod              *corev1.Pod
}

func (b *Build) Init(args map[string]interface{}) error {
  // Decode/assign args to struct keys.
  b.ResourceID = args["resourceID"].(uint)
  b.BuildTargetSha = args["buildTargetSha"].(string)
  b.TargetCluster = args["targetCluster"].(string)
  b.LogStreamKey = args["logKey"].(string)

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

  // Initialize the result channel.
  b.ResultChannel = make(chan Result)

  // Set project and builable resource.
  b.Project = project
  b.Buildable = resource

  // Cluster name is the SweetTea Build Cluster name.
  b.ClusterName = app.Config.BuildClusterName

  // Create unique container and deploy names.
  b.ContainerName = fmt.Sprintf("%s-%s-%s", b.TargetCluster, cluster.Build, project.Uid)
  b.DeployName = fmt.Sprintf("%s-%v", b.ContainerName, timeutil.MSSinceEpoch())

  // Set image to the SweetTea Build Server image.
  b.Image = fmt.Sprintf("%s/%s:%s", app.Config.DockerRegistryOrg, image.BuildServer, app.Config.BuildServerImageVersion)

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

// Perform deploys the configured pod to the Build Cluster.
func (b *Build) Perform() error {
  return CreatePod(b.Client, b.Namespace, b.Pod, cluster.Build)
}

func (b *Build) GetResultChannel() chan Result {
  return b.ResultChannel
}

func (b *Build) Watch() {
  // Get a namespaced pod-watcher channel.
  ch, err := PodWatcherChannel(b.Client, b.Namespace, b.DeployName)

  if err != nil {
    b.ResultChannel <- Result{Ok: false, Error: err}
    return
  }

  // Start watching for events.
  for event := range ch {
    if result := b.checkEventForResult(event); result != nil {
      b.ResultChannel <- *result
      return
    }
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
    "BUILD_TARGET_ACCESS_TOKEN": projectsvc.GetHost(b.Project).GetToken(),
    "BUILD_TARGET_SHA": b.BuildTargetSha,
    "BUILD_TARGET_UID": b.Project.Uid,
    "BUILD_TARGET_URL": b.Project.Url(),
    "LOG_STREAM_KEY": b.LogStreamKey,
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
    "name": b.DeployName,
    "containers": b.Containers,
    "volumes": b.Volumes,
    "restart": corev1.RestartPolicyNever,
  })
}

func (b *Build) checkEventForResult(event watch.Event) *Result {
  switch event.Type {

  case watch.Added:
    // Log when pod is added.
    app.Log.Infof("Job %s started.", b.DeployName)
    return nil

  case watch.Modified:
    // When pod is modified, check its status and report result when success or failure.
    pod, ok := event.Object.(*corev1.Pod)

    if !ok {
      err := fmt.Errorf("Job %s encountered unexpected event object type.", b.DeployName)
      app.Log.Errorf(err.Error())
      return &Result{Ok: false, Error: err}
    }
    
    // Check for pod success/failure.
    return b.checkPodStatus(pod.Status)

  case watch.Error:
    // Error out during a pod error.
    err := fmt.Errorf("Job %s encountered pod error.", b.DeployName)
    app.Log.Errorf(err.Error())
    return &Result{Ok: false, Error: err}

  default:
    return nil
  }
}

func (b *Build) checkPodStatus(podStatus corev1.PodStatus) *Result {
  switch podStatus.Phase {

  case corev1.PodSucceeded:
    app.Log.Infoln("Successfully built image.")
    return &Result{Ok: true}

  case corev1.PodFailed:
    err := fmt.Errorf("Job %s failed with error: -- %s.", b.DeployName, podStatus.Message)
    app.Log.Errorf(err.Error())
    return &Result{Ok: false, Error: err}

  case corev1.PodUnknown:
    err := fmt.Errorf("Job %s encountered unknown pod status error: %s.", b.DeployName, podStatus.Message)
    app.Log.Errorf(err.Error())
    return &Result{Ok: false, Error: err}

  default:
    return nil
  }
}