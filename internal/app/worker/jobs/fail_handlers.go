package jobs

import (
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/pkg/service/buildablesvc"
  "github.com/sweettea-io/rest-api/internal/pkg/service/deploysvc"
  "github.com/sweettea-io/rest-api/internal/pkg/service/trainjobsvc"
)

func failBuildable(buildableID uint, targetCluster string, err error, logKey string, logMsg string, args ...interface{}) error {
  if e := buildablesvc.Fail(buildableID, targetCluster); e != nil {
    app.Log.Errorln(e.Error())
  }

  return logBuildableErr(err, logKey, logMsg, args...)

}

func failDeploy(deployID uint, err error, deployLogKey string, deployLogMsg string, args ...interface{}) error {
  if e := deploysvc.FailByID(deployID); e != nil {
    app.Log.Errorln(e.Error())
  }

  return logBuildableErr(err, deployLogKey, deployLogMsg, args...)
}

func failTrainJob(trainJobID uint, err error, trainJobLogKey string, trainJobLogMsg string, args ...interface{}) error {
  if e := trainjobsvc.FailByID(trainJobID); e != nil {
    app.Log.Errorln(e.Error())
  }

  return logBuildableErr(err, trainJobLogKey, trainJobLogMsg, args...)
}

func logBuildableErr(err error, buildableLogKey string, buildableLogMsg string, args ...interface{}) error {
  app.Log.Errorln(err.Error())
  app.Log.BuildableErrorf(buildableLogKey, buildableLogMsg, args...)
  return err
}