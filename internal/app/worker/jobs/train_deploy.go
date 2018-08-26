package jobs

import (
  "fmt"
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/pkg/service/trainjobsvc"
  "github.com/sweettea-io/work"
  "github.com/sweettea-io/rest-api/internal/pkg/service/envvarsvc"
  "github.com/sweettea-io/rest-api/internal/pkg/k"
  "github.com/sweettea-io/rest-api/internal/pkg/model"
)

/*
  TrainDeploy deploys a TrainJob to the SweetTea Train Cluster for model training.

  Args:
    trainJobID (uint)   ID of the TrainJob to deploy
    envs       (string) json string representation of the custom env vars to use with this Train Cluster deploy
    logKey     (string) log key for buildable
*/
func (c *Context) TrainDeploy(job *work.Job) error {
  // Extract args from job.
  trainJobID := uint(job.ArgInt64("trainJobID"))
  envs := job.ArgString("envs")
  logKey := job.ArgString("logKey")

  if err := job.ArgError(); err != nil {
    if logKey != "" {
      return failTrainJob(trainJobID, err, logKey, "Arg error occurred inside train deploy job.")
    }

    app.Log.Errorln(err.Error())
    return err
  }

  // Ensure Train Cluster exists first.
  if !app.Config.TrainClusterConfigured() {
    err := fmt.Errorf("Train Cluster not yet configured.")
    return failTrainJob(trainJobID, err, logKey, err.Error())
  }

  // Convert stringified envs into map[string]string representation.
  envsMap, err := envvarsvc.MapFromBytes([]byte(envs))
  if err != nil {
    return failTrainJob(trainJobID, err, logKey, "Failed to parse train environment variables.")
  }

  // Create K8S train deploy and prep args.
  trainDeploy := k.Train{}
  trainDeployArgs := map[string]interface{}{
    "trainJobID": trainJobID,
    "envs": envsMap,
  }

  // Initialize train deploy.
  if err := trainDeploy.Init(trainDeployArgs); err != nil {
    return failTrainJob(trainJobID, err, logKey, "Failed to initialize train deploy.")
  }

  // Create deploy resources.
  if err := trainDeploy.Configure(); err != nil {
    return failTrainJob(trainJobID, err, logKey, "Failed to configure train deploy resources.")
  }

  // Deploy to train cluster.
  if err := trainDeploy.Perform(); err != nil {
    return failTrainJob(trainJobID, err, logKey, "Failed to perform train deploy.")
  }

  // Update TrainJob stage to Deployed.
  if err := trainjobsvc.UpdateStageByID(trainJobID, model.BuildStages.Deployed); err != nil {
    return failTrainJob(trainJobID, err, logKey, "Failed to update stage of train job.")
  }

  // Get channel to watch train deploy.
  resultCh := trainDeploy.GetResultChannel()

  // Watch train deploy until the pod has successfully been added.
  go trainDeploy.Watch()
  deployResult := <-resultCh

  // Error out if pod failed to be added.
  if !deployResult.Ok {
    return failTrainJob(trainJobID, deployResult.Error, logKey, "Deploy to Train Cluster failed.")
  }

  // TODO: Stream message back successfully disconnecting client.

  return nil
}