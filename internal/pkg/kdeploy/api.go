package kdeploy

import (
  "fmt"
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/pkg/model"
  "github.com/sweettea-io/rest-api/internal/pkg/model/buildable"
  "github.com/sweettea-io/rest-api/internal/pkg/service/deploysvc"
  "github.com/sweettea-io/rest-api/internal/pkg/service/envvarsvc"
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
  Project         *model.Project
  ModelVersion    *model.ModelVersion
  Model           *model.Model
  ApiCluster      *model.ApiCluster
  DeployName      string
  Image           string
  ContainerName   string
  ResultChannel   <-chan Result

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
  // Find Deploy by ID.
  deploy, err := deploysvc.FromID(args["deployID"].(uint))

  if err != nil {
    return err
  }

  // Update Deploy to Deploying.
  if err := deploysvc.UpdateStage(deploy, buildable.Deploying); err != nil {
    return err
  }

  // Initialize the result channel.
  api.ResultChannel = make(chan Result)

  // Store refs to models.
  api.Deploy = deploy
  api.Commit = &deploy.Commit
  api.Project = &deploy.Commit.Project
  api.ModelVersion = &deploy.ModelVersion
  api.Model = &deploy.ModelVersion.Model
  api.ApiCluster = &deploy.ApiCluster

  // Name of container to be run inside the pods.
  api.ContainerName = fmt.Sprintf("%s-%s", cluster.Api, api.Project.Uid)

  // Docker image to deploy (ex: sweetteaprod/api-<project_uid>:<commit_sha>)
  api.Image = fmt.Sprintf("%s/%s:%s", app.Config.DockerRegistryOrg, api.ContainerName, api.Commit.Sha)

  // Create a new unique deploy name for new deploys. Otherwise, use the existing deploy name.
  if api.IsNewDeployment() {
    api.DeployName = fmt.Sprintf("%s-%v", api.ContainerName, timeutil.MSSinceEpoch())
  } else {
    api.DeployName = api.Deploy.Name
  }

  return nil
}

func (api *Api) Configure() error {
  // Configure k8s client.
  if err := api.makeClient(); err != nil {
    return err
  }

  // Build up k8s resources into the Pod we wish to deploy.
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
  // What patching an existing deploy would look like:
  //jsonReprOfStuffToChange := `{shit to change}`
  //api.Client.Deployments(api.Namespace).Patch(api.DeployName, types.JSONPatchType, []byte(jsonReprOfStuffToChange))

  return CreateDeployment(api.Client, api.Namespace, api.Deployment)
}

func (api *Api) GetResultChannel() <-chan Result {
  return api.ResultChannel
}

func (api *Api) Watch() {
  // Get a namespaced deployment-watcher channel.
  ch, err := DeploymentWatcherChannel(api.Client, api.Namespace, api.DeployName)

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

// IsNewDeployment returns whether a k8s deployment already exists for this Deploy.
func (api *Api) IsNewDeployment() bool {
  return api.Deploy.Name == ""
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
  // Start with Deploy EnvVars.
  envs := envvarsvc.GetMap(api.Deploy.ID)

  // Add deploy-agnostic envs.
  envs = maputil.MergeMaps(envs, map[string]string{
    "AWS_ACCESS_KEY_ID": app.Config.AWSAccessKeyId,
    "AWS_REGION_NAME": app.Config.AWSRegionName,
    "AWS_SECRET_ACCESS_KEY": app.Config.AWSSecretAccessKey,
    "LOG_STREAM_KEY": api.Deploy.Uid,
    "MODEL_STORAGE_URL": app.Config.ModelStorageUrl,
    "MODEL_KEY": api.ModelVersion.StorageKey(),
    "PROJECT_UID": api.Project.Uid,
    "REDIS_ADDRESS": app.Config.RedisAddress,
    "REDIS_PASSWORD": app.Config.RedisPassword,
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
    "label": api.DeployName,
    "containers": api.Containers,
    "restart": corev1.RestartPolicyNever,
  })
}

func (api *Api) makeDeploymentSpec() {
  api.DeploymentSpec = DeploymentSpec(api.PodTemplateSpec, app.Config.DeployReplicasCount)
}

func (api *Api) makeDeployment() {
  api.Deployment = Deployment(api.DeploymentSpec, api.DeployName)
}

func (api *Api) checkEventForResult(event watch.Event) *Result {
  switch event.Type {

  // Log & return with success when deployment has been added.
  case watch.Added:
    app.Log.Infof("Deployment %s started.", api.DeployName)
    return &Result{Ok: true}

  // Return with error if watch error occurs before deployment has been added.
  case watch.Error:
    err := fmt.Errorf("Deployment %s encountered error.", api.DeployName)
    app.Log.Errorf(err.Error())
    return &Result{Ok: false, Error: err}

  default:
    return nil
  }
}