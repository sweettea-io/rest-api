package jobs

import (
  "fmt"
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/pkg/model/buildable"
  "github.com/sweettea-io/rest-api/internal/pkg/service/trainjobsvc"
  "github.com/sweettea-io/work"
  "github.com/sweettea-io/rest-api/internal/pkg/service/envvarsvc"
  "github.com/sweettea-io/rest-api/internal/pkg/k"
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
    return failTrainJob(trainJobID, err)
  }

  // Ensure Train Cluster exists first.
  if !app.Config.TrainClusterConfigured() {
    return failTrainJob(trainJobID, fmt.Errorf("train cluster not configured -- leaving CreateTrainJob"))
  }

  // Convert stringified envs into map[string]string representation.
  envsMap, err := envvarsvc.MapFromBytes([]byte(envs))

  if err != nil {
    return failTrainJob(trainJobID, err)
  }

  // Create K8S train deploy and prep args.
  trainDeploy := k.Train{}
  trainDeployArgs := map[string]interface{}{
    "trainJobID": trainJobID,
    "envs": envsMap,
  }

  // Initialize train deploy.
  if err := trainDeploy.Init(trainDeployArgs); err != nil {
    return failTrainJob(trainJobID, err)
  }

  // Create deploy resources.
  if err := trainDeploy.Configure(); err != nil {
    return failTrainJob(trainJobID, err)
  }

  // Deploy to train cluster.
  if err := trainDeploy.Perform(); err != nil {
    return failTrainJob(trainJobID, err)
  }

  // Update TrainJob stage to Deployed.
  if err := trainjobsvc.UpdateStageByID(trainJobID, buildable.Deployed); err != nil {
    return failTrainJob(trainJobID, err)
  }

  // TODO: Stream message back successfully disconnecting client.

  // Get channel to watch train deploy.
  resultCh := trainDeploy.GetResultChannel()

  // Watch train deploy until the pod has successfully been added.
  go trainDeploy.Watch()
  deployResult := <-resultCh

  // Error out if pod failed to be added.
  if !deployResult.Ok {
    return failTrainJob(trainJobID, deployResult.Error)
  }

  return nil
}