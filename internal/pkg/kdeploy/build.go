package kdeploy

import (
  "fmt"
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/pkg/model"
  "github.com/sweettea-io/rest-api/internal/pkg/service/projectsvc"
  "github.com/sweettea-io/rest-api/internal/pkg/util/cluster"
  corev1 "k8s.io/api/core/v1"
  typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type Build struct {
  // Args
  ResourceID    uint
  TargetCluster string

  // Created on Init
  Project       *model.Project
  ClusterName   string

  // K8S resources
  Namespace     string
  Client        *typedcorev1.CoreV1Client
  VolumeMounts  []*corev1.VolumeMount
  Volumes       []*corev1.Volume
  Envs          []*corev1.EnvVar
  Container     *corev1.Container
  Pod           *corev1.Pod
}

func (b *Build) Init(args map[string]interface{}) error {
  // Decode/assign args to struct keys.
  b.ResourceID = args["resourceID"].(uint)
  b.TargetCluster = args["targetCluster"].(string)
  b.ClusterName = app.Config.BuildClusterName

  // Find Project by ID.
  project, err := projectsvc.FromID(args["projectID"].(uint))

  if err != nil {
    return err
  }

  b.Project = project
  return nil
}

func (b *Build) Configure() error {
  steps := []func() error {
    b.makeClient,
    b.makeVolumeMounts,
    b.makeVolumes,
    b.makeEnvs,
    b.makeContainer,
    b.makePod,
  }

  for _, step := range steps {
    if err := step(); err != nil {
      return err
    }
  }

  return nil
}

func (b *Build) Perform() error {
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

func (b *Build) makeVolumeMounts() error {
  return nil
}

func (b *Build) makeVolumes() error {
  return nil
}

func (b *Build) makeEnvs() error {
  return nil
}

func (b *Build) makeContainer() error {
  return nil
}

func (b *Build) makePod() error {
  return nil
}