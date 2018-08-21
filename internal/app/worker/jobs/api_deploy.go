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
    deployID (uint) ID of the Deploy model
*/
func (c *Context) ApiDeploy(job *work.Job) error {
  // Extract args from job.
  deployID := uint(job.ArgInt64("deployID"))

  if err := job.ArgError(); err != nil {
    deploysvc.FailByID(deployID)
    app.Log.Errorln(err.Error())
    return err
  }

  // Get Deploy by ID.
  deploy, err := deploysvc.FromID(deployID)

  if err != nil {
    deploysvc.FailByID(deployID)
    app.Log.Errorln(err.Error())
    return err
  }

  // Create K8S API deploy and prep args.
  apiDeploy := k.Api{}

  apiDeployArgs := map[string]interface{}{
    "deploy": deploy,
    "commit": &deploy.Commit,
    "model": &deploy.Model,
    "customEnvs": envvarsvc.GetMap(deploy.ID),
  }

  // Initialize API deploy.
  if err := apiDeploy.Init(apiDeployArgs); err != nil {
    deploysvc.FailByID(deployID)
    app.Log.Errorln(err.Error())
    return err
  }

  // Create deploy resources.
  if err := apiDeploy.Configure(); err != nil {
    deploysvc.FailByID(deployID)
    app.Log.Errorln(err.Error())
    return err
  }

  // Deploy to ApiCluster.
  if err := apiDeploy.Perform(); err != nil {
    deploysvc.FailByID(deployID)
    app.Log.Errorln(err.Error())
    return err
  }

  // Update Deploy stage to Deployed and register its deployment name.
  updates := map[string]interface{}{
    "deployment_name": apiDeploy.DeploymentName,
  }

  if err := deploysvc.Deployed(deployID, updates); err != nil {
    deploysvc.FailByID(deployID)
    app.Log.Errorln(err.Error())
    return err
  }

  // TODO: Stream message back successfully disconnecting client.

  // Get channel to watch API deploy.
  resultCh := apiDeploy.GetResultChannel()

  // Watch API deployment until it has successfully started.
  go apiDeploy.Watch()
  deployResult := <-resultCh

  // Error out if deployment failed to start.
  if !deployResult.Ok {
    deploysvc.FailByID(deployID)
    app.Log.Errorf(deployResult.Error.Error())
    return deployResult.Error
  }

  // Schedule Deploy publication.
  if _, err := app.JobQueue.Enqueue(Names.PublicizeDeploy, work.Q{"deployID": deployID}); err != nil {
    deploysvc.FailByID(deployID)
    app.Log.Errorf("error scheduling PublicizeDeploy job: %s", err.Error())
    return err
  }

  return nil
}