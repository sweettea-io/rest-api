package k

import (
  "fmt"
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/pkg/model"
  "github.com/sweettea-io/rest-api/internal/pkg/service/deploysvc"
  "github.com/sweettea-io/rest-api/internal/pkg/util/cluster"
  "github.com/sweettea-io/rest-api/internal/pkg/util/maputil"
  "github.com/sweettea-io/rest-api/internal/pkg/util/timeutil"
  "k8s.io/api/extensions/v1beta1"
  "k8s.io/apimachinery/pkg/watch"
  corev1 "k8s.io/api/core/v1"
  typedbetav1 "k8s.io/client-go/kubernetes/typed/extensions/v1beta1"
)

type Api struct {
  // Establish on Init
  Deploy          *model.Deploy
  Commit          *model.Commit
  ModelVersion    *model.ModelVersion
  CustomEnvs      map[string]string
  LogStreamKey    string
  Project         *model.Project
  Model           *model.Model
  ApiCluster      *model.ApiCluster
  DeploymentName  string
  Image           string
  ContainerName   string
  ResultChannel   chan Result

  // K8S resources
  Namespace       string
  Client          *typedbetav1.ExtensionsV1beta1Client
  Envs            []corev1.EnvVar
  Ports           []corev1.ContainerPort
  Containers      []corev1.Container
  PodTemplateSpec *corev1.PodTemplateSpec
  DeploymentSpec  *v1beta1.DeploymentSpec
  Deployment      *v1beta1.Deployment
}

func (api *Api) Init(args map[string]interface{}) error {
  // Parse args.
  api.Deploy = args["deploy"].(*model.Deploy)
  api.Commit = args["commit"].(*model.Commit)
  api.ModelVersion = args["modelVersion"].(*model.ModelVersion)
  api.CustomEnvs = args["customEnvs"].(map[string]string)
  api.LogStreamKey = args["logKey"].(string)

  // Update Deploy to Deploying.
  if err := deploysvc.UpdateStage(api.Deploy, model.BuildStages.Deploying); err != nil {
    return err
  }

  // Initialize the result channel.
  api.ResultChannel = make(chan Result)

  // Set further models through associations.
  api.Project = &api.Commit.Project
  api.Model = &api.ModelVersion.Model
  api.ApiCluster = &api.Deploy.ApiCluster

  // Name of container to be run inside the pods.
  api.ContainerName = fmt.Sprintf("%s-%s", cluster.Api, api.Project.Uid)

  // Docker image to deploy (ex: sweetteaprod/api-<project_uid>:<commit_sha>)
  api.Image = fmt.Sprintf("%s/%s:%s", app.Config.DockerRegistryOrg, api.ContainerName, api.Commit.Sha)

  if api.IsNewDeployment() {
    // Create new deployment name if not.
    api.DeploymentName = fmt.Sprintf("%s-%v", api.ContainerName, timeutil.MSSinceEpoch())
  } else {
    // Use current deployment if already exists.
    api.DeploymentName = api.Deploy.DeploymentName
  }

  return nil
}

func (api *Api) Configure() error {
  // Configure k8s client.
  if err := api.makeClient(); err != nil {
    return err
  }

  // Build up K8S resources into the Pod we wish to deploy.
  api.makeEnvs()
  api.makePorts()
  api.makeContainers()
  api.makePodTemplateSpec()
  api.makeDeploymentSpec()
  api.makeDeployment()

  return nil
}

// Perform deploys the configured deployment to provided ApiCluster.
func (api *Api) Perform() error {
  if api.IsNewDeployment() {
    return CreateDeployment(api.Client, api.Namespace, api.Deployment)
  } else {
    return UpdateDeployment(api.Client, api.Namespace, api.Deployment)
  }
}

func (api *Api) GetResultChannel() chan Result {
  return api.ResultChannel
}

func (api *Api) Watch() {
  // Get a namespaced deployment watcher channel.
  ch, err := DeploymentWatcherChannel(api.Client, api.Namespace, api.DeploymentName)

  if err != nil {
    api.ResultChannel <- Result{Ok: false, Error: err}
    return
  }

  // Start watching for events.
  for event := range ch {
    if result := api.checkEventForResult(event); result != nil {
      api.ResultChannel <- *result
      return
    }
  }
}

func (api *Api) IsNewDeployment() bool {
  return api.Deploy.DeploymentName == ""
}

func (api *Api) makeClient() error {
  // Configure V1Beta1 client.
  client, nsp, err := ConfigureV1Beta1(api.ApiCluster.Slug)

  if err != nil {
    return err
  }

  // Store refs to client and namespace.
  api.Client = client
  api.Namespace = nsp

  return nil
}

func (api *Api) makeEnvs() {
  envs := maputil.MergeMaps(api.CustomEnvs, map[string]string{
    "AWS_ACCESS_KEY_ID": app.Config.AWSAccessKeyId,
    "AWS_REGION_NAME": app.Config.AWSRegionName,
    "AWS_SECRET_ACCESS_KEY": app.Config.AWSSecretAccessKey,
    "CLIENT_ID": api.Deploy.ClientID,
    "CLIENT_SECRET": api.Deploy.ClientSecret,
    "LOG_STREAM_KEY": api.LogStreamKey,
    "MODEL_STORAGE_URL": app.Config.ModelStorageUrl,
    "MODEL_STORAGE_FILE_PATH": api.ModelVersion.StorageKey(api.Project, api.Model),
    "PROJECT_UID": api.Project.Uid,
  })

  // Format envs as k8s resources.
  api.Envs = EnvVars(envs)
}

func (api *Api) makePorts() {
  api.Ports = Ports([]int32{80})
}

func (api *Api) makeContainers() {
  api.Containers = Containers([]map[string]interface{}{{
    "name": api.ContainerName,
    "image": api.Image,
    "envs": api.Envs,
    "ports": api.Ports,
  }})
}

func (api *Api) makePodTemplateSpec() {
  api.PodTemplateSpec = PodTemplateSpec(map[string]interface{}{
    "name": api.DeploymentName,
    "containers": api.Containers,
    "restart": corev1.RestartPolicyNever,
  })
}

func (api *Api) makeDeploymentSpec() {
  api.DeploymentSpec = DeploymentSpec(api.PodTemplateSpec, app.Config.DeployReplicasCount)
}

func (api *Api) makeDeployment() {
  api.Deployment = Deployment(api.DeploymentSpec, api.DeploymentName)
}

func (api *Api) checkEventForResult(event watch.Event) *Result {
  switch event.Type {

  // Log & return with success when deployment has been added.
  case watch.Added:
    app.Log.Infof("Deployment %s started.", api.DeploymentName)
    return &Result{Ok: true}

  // Return with error if watch error occurs before deployment has been added.
  case watch.Error:
    err := fmt.Errorf("Deployment %s encountered error.", api.DeploymentName)
    app.Log.Errorf(err.Error())
    return &Result{Ok: false, Error: err}

  default:
    return nil
  }
}