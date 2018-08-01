package jobs

import (
  "fmt"
  "encoding/json"
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/work"
  "github.com/sweettea-io/rest-api/internal/pkg/kdeploy"
  "github.com/sweettea-io/rest-api/internal/pkg/service/trainjobsvc"
  "github.com/sweettea-io/rest-api/internal/pkg/model/buildable"
)

/*
  TrainDeploy deploys a TrainJob to the SweetTea Train Cluster for model training.

  Args:
    trainJobID (uint) ID of the TrainJob to deploy
    envs       (string) json string representation of the custom env vars to use with this Train Cluster deploy
*/
func (c *Context) TrainDeploy(job *work.Job) error {
  // Extract args from job.
  trainJobID := uint(job.ArgInt64("trainJobID"))
  envs := job.ArgString("envs")

  if err := job.ArgError(); err != nil {
    trainjobsvc.FailByID(trainJobID)
    app.Log.Errorln(err.Error())
    return err
  }

  // Ensure Train Cluster exists first.
  if !app.Config.TrainClusterConfigured() {
    err := fmt.Errorf("train cluster not configured -- leaving CreateTrainJob")
    trainjobsvc.FailByID(trainJobID)
    app.Log.Errorln(err.Error())
    return err
  }

  // Convert stringified envs into map[string]string representation.
  var envsMap map[string]string
  if err := json.Unmarshal([]byte(envs), &envsMap); err != nil {
    err = fmt.Errorf("error converting custom train envs into map[string]string: %s", err.Error())
    trainjobsvc.FailByID(trainJobID)
    app.Log.Errorln(err.Error())
    return err
  }

  // Create K8S train deploy and prep args.
  trainDeploy := kdeploy.Train{}

  trainDeployArgs := map[string]interface{}{
    "trainJobID": trainJobID,
    "envs": envsMap,
  }

  // Initialize train deploy.
  if err := trainDeploy.Init(trainDeployArgs); err != nil {
    trainjobsvc.FailByID(trainJobID)
    app.Log.Errorln(err.Error())
    return err
  }

  // Create deploy resources.
  if err := trainDeploy.Configure(); err != nil {
    trainjobsvc.FailByID(trainJobID)
    app.Log.Errorln(err.Error())
    return err
  }

  // Deploy to train cluster.
  if err := trainDeploy.Perform(); err != nil {
    trainjobsvc.FailByID(trainJobID)
    app.Log.Errorln(err.Error())
    return err
  }

  // Update TrainJob stage to Deployed.
  if err := trainjobsvc.UpdateStageByID(trainJobID, buildable.Deployed); err != nil {
    trainjobsvc.FailByID(trainJobID)
    app.Log.Errorln(err.Error())
    return err
  }

  return nil
}