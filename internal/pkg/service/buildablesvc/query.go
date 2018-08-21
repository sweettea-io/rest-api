package buildablesvc

import (
  "fmt"
  "github.com/sweettea-io/rest-api/internal/pkg/service/deploysvc"
  "github.com/sweettea-io/rest-api/internal/pkg/service/trainjobsvc"
  "github.com/sweettea-io/rest-api/internal/pkg/util/cluster"
  "github.com/sweettea-io/rest-api/internal/pkg/model"
)

func FromID(id uint, targetCluster string) (model.Buildable, error) {
  switch targetCluster {
  case cluster.Train:
    return trainjobsvc.FromID(id)
  case cluster.Api:
    return deploysvc.FromID(id)
  default:
    return nil, fmt.Errorf("couldn't find buildable for target cluster %s", targetCluster)
  }
}