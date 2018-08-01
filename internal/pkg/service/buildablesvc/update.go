package buildablesvc

import (
  "fmt"
  "github.com/sweettea-io/rest-api/internal/pkg/util/cluster"
  "github.com/sweettea-io/rest-api/internal/pkg/service/trainjobsvc"
  "github.com/sweettea-io/rest-api/internal/pkg/service/deploysvc"
)

func UpdateStage(buildableID uint, stage string, targetCluster string) error {
  switch targetCluster {
  case cluster.Train:
    return trainjobsvc.UpdateStageByID(buildableID, stage)
  case cluster.Api:
    return deploysvc.UpdateStageByID(buildableID, stage)
  default:
    return fmt.Errorf("couldn't update buildable(id=%v) stage for target cluster %s", buildableID, targetCluster)
  }
}

func Fail(buildableID uint, targetCluster string) error {
  switch targetCluster {
  case cluster.Train:
    return trainjobsvc.FailByID(buildableID)
  case cluster.Api:
    return deploysvc.FailByID(buildableID)
  default:
    return fmt.Errorf("couldn't update buildable(id=%v) stage for target cluster %s", buildableID, targetCluster)
  }
}