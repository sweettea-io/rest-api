package jobs

import (
  "fmt"
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/work"
  "github.com/sweettea-io/rest-api/internal/pkg/kdeploy"
  "encoding/json"
)

/*
  TrainDeploy deploys a TrainJob to the SweetTea Train Cluster for model training.

  Args:
    trainJobID (uint) ID of the TrainJob to deploy
    envs       (string) json string representation of the custom env vars to use with this Train Cluster deploy
*/
func (c *Context) TrainDeploy(job *work.Job) error {
  // Ensure Train cluster exists first.
  if !app.Config.TrainClusterConfigured() {
    err := fmt.Errorf("train cluster not configured -- leaving CreateTrainJob")
    app.Log.Errorln(err.Error())
    return err
  }

  // Extract args from job.
  trainJobID := uint(job.ArgInt64("trainJobID"))
  envs := job.ArgString("envs")

  if err := job.ArgError(); err != nil {
    app.Log.Errorln(err.Error())
    return err
  }

  // Convert stringified envs into map[string]string representation.
  var envsMap map[string]string
  if err := json.Unmarshal([]byte(envs), &envsMap); err != nil {
    err = fmt.Errorf("error converting custom train envs into map[string]string: %s", err.Error())
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
    app.Log.Errorln(err.Error())
    return err
  }

  // Create deploy resources.
  if err := trainDeploy.Configure(); err != nil {
    app.Log.Errorln(err.Error())
    return err
  }

  // Deploy to train cluster.
  if err := trainDeploy.Perform(); err != nil {
    app.Log.Errorln(err.Error())
    return err
  }

  return nil
}