package jobs

import (
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/pkg/service/buildablesvc"
  "github.com/sweettea-io/rest-api/internal/pkg/service/deploysvc"
  "github.com/sweettea-io/rest-api/internal/pkg/service/trainjobsvc"
)

func failBuildable(buildableID uint, targetCluster string, err error) error {
  if e := buildablesvc.Fail(buildableID, targetCluster); e != nil {
    return logAndFail(e)
  }

  return logAndFail(err)

}

func failDeploy(deployID uint, err error) error {
  if e := deploysvc.FailByID(deployID); e != nil {
    return logAndFail(e)
  }

  return logAndFail(err)
}

func failTrainJob(trainJobID uint, err error) error {
  if e := trainjobsvc.FailByID(trainJobID); e != nil {
    return logAndFail(e)
  }

  return logAndFail(err)
}

func logAndFail(err error) error {
  app.Log.Errorln(err.Error())
  return err
}