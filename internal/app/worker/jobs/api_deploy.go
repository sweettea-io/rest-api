package jobs

import (
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/pkg/service/deploysvc"
  "github.com/sweettea-io/work"
  "github.com/sweettea-io/rest-api/internal/pkg/service/envvarsvc"
  "github.com/sweettea-io/rest-api/internal/pkg/k"
)

/*
  ApiDeploy deploys a Deploy model to an ApiCluster.

  Args:
    deployID (uint)   ID of the Deploy model
    logKey   (string) log key for buildable
*/
func (c *Context) ApiDeploy(job *work.Job) error {
  // Extract args from job.
  deployID := uint(job.ArgInt64("deployID"))
  logKey := job.ArgString("logKey")

  if err := job.ArgError(); err != nil {
    if logKey != "" {
      return failDeploy(deployID, err, logKey, "Arg error occurred inside API deploy job.")
    }

    app.Log.Errorln(err.Error())
    return err
  }

  // Get Deploy by ID.
  deploy, err := deploysvc.FromID(deployID)
  if err != nil {
    return failDeploy(deployID, err, logKey, "Deploy not found.")
  }

  // Create K8S API deploy and prep args.
  apiDeploy := k.Api{}

  apiDeployArgs := map[string]interface{}{
    "deploy": deploy,
    "commit": &deploy.Commit,
    "model": &deploy.Model,
    "customEnvs": envvarsvc.GetMap(deploy.ID),
    "logKey": logKey,
  }

  // Initialize API deploy.
  if err := apiDeploy.Init(apiDeployArgs); err != nil {
    return failDeploy(deployID, err, logKey, "Failed to initialize API deploy.")
  }

  // Create deploy resources.
  if err := apiDeploy.Configure(); err != nil {
    return failDeploy(deployID, err, logKey, "Failed to configure API deploy resources.")
  }

  // Deploy to ApiCluster.
  if err := apiDeploy.Perform(); err != nil {
    return failDeploy(deployID, err, logKey, "Failed to perform API deploy.")
  }

  // Update Deploy stage to Deployed and register its deployment name.
  updates := map[string]interface{}{
    "deployment_name": apiDeploy.DeploymentName,
  }

  if err := deploysvc.Deployed(deployID, updates); err != nil {
    return failDeploy(deployID, err, logKey, "Failed to update stage of deploy.")
  }

  // Get channel to watch API deploy.
  resultCh := apiDeploy.GetResultChannel()

  // Watch API deployment until it has successfully started.
  go apiDeploy.Watch()
  deployResult := <-resultCh

  // Error out if deployment failed to start.
  if !deployResult.Ok {
    return failDeploy(deployID, err, logKey, "Deploy to API cluster failed.")
  }

  // Schedule Deploy publication.
  if _, err := app.JobQueue.Enqueue(Names.PublicizeDeploy, work.Q{"deployID": deployID, "logKey": logKey}); err != nil {
    return failDeploy(deployID, err, logKey, "Failed to schedule publicize deploy job.")
  }

  return nil
}