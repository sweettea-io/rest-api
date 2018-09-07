package k

import (
  "fmt"
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/pkg/model"
  "github.com/sweettea-io/rest-api/internal/pkg/service/trainjobsvc"
  "github.com/sweettea-io/rest-api/internal/pkg/util/cluster"
  "github.com/sweettea-io/rest-api/internal/pkg/util/maputil"
  "github.com/sweettea-io/rest-api/internal/pkg/util/timeutil"
  "k8s.io/apimachinery/pkg/watch"
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
  ResultChannel chan Result

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

  // Update TrainJob to Deploying.
  if err := trainjobsvc.UpdateStage(trainJob, model.BuildStages.Deploying); err != nil {
    return err
  }

  // Initialize the result channel.
  t.ResultChannel = make(chan Result)

  // Store refs to models.
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
  return CreatePod(t.Client, t.Namespace, t.Pod, cluster.Train)
}


func (t *Train) GetResultChannel() chan Result {
  return t.ResultChannel
}

func (t *Train) Watch() {
  // Get a namespaced pod-watcher channel.
  ch, err := PodWatcherChannel(t.Client, t.Namespace, t.DeployName)

  if err != nil {
    t.ResultChannel <- Result{Ok: false, Error: err}
    return
  }

  // Start watching for events.
  for event := range ch {
    if result := t.checkEventForResult(event); result != nil {
      t.ResultChannel <- *result
      return
    }
  }
}

func (t *Train) makeClient() error {
  // Configure CoreV1 client.
  client, nsp, err := ConfigureCoreV1(t.ClusterName)

  if err != nil {
    return err
  }

  // Store refs to client and namespace.
  t.Client = client
  t.Namespace = nsp

  return nil
}

func (t *Train) makeEnvs() {
  // Create envs map by merging deploy-agnostic envs on top of custom envs.
  envs := maputil.MergeMaps(t.CustomEnvs, map[string]string{
    "AWS_ACCESS_KEY_ID": app.Config.AWSAccessKeyId,
    "AWS_REGION_NAME": app.Config.AWSRegionName,
    "AWS_SECRET_ACCESS_KEY": app.Config.AWSSecretAccessKey,
    "LOG_STREAM_KEY": t.TrainJob.Uid,
    "MODEL_STORAGE_URL": app.Config.ModelStorageUrl,
    "MODEL_STORAGE_FILE_PATH": t.ModelVersion.StorageKey(t.Project, t.Model),
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
    "name": t.DeployName,
    "containers": t.Containers,
    "restart": corev1.RestartPolicyNever,
  })
}

func (t *Train) checkEventForResult(event watch.Event) *Result {
  switch event.Type {

  // Log & return with success when Training pod has been added.
  case watch.Added:
    app.Log.Infof("Job %s started.", t.DeployName)
    return &Result{Ok: true}

  // Return with error if watch error occurs before pod has been added.
  case watch.Error:
    err := fmt.Errorf("Job %s encountered pod error.", t.DeployName)
    app.Log.Errorf(err.Error())
    return &Result{Ok: false, Error: err}

  default:
    return nil
  }
}