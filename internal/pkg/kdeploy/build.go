package kdeploy

import (
  "github.com/sweettea-io/rest-api/internal/pkg/service/projectsvc"
  "github.com/sweettea-io/rest-api/internal/pkg/model"
  "github.com/sweettea-io/rest-api/internal/pkg/util/cluster"
  "github.com/sweettea-io/rest-api/internal/app"
)

type Build struct {
  ResourceID    uint
  Project       *model.Project
  TargetCluster string
  ClusterName   string
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
  // Get new rest client config from kubeconfig and desired context.
  restConfig, nsp, err := ConfigureClient(b.ClusterName)

  if err != nil {
    return err
  }

  // Get new CoreV1 client.
  coreV1Client, err := NewCoreV1(restConfig)

  if err != nil {
    return err
  }

  return nil
}

func (b *Build) Perform() error {
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