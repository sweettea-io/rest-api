package jobs

import (
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/pkg/kdeploy"
  "github.com/sweettea-io/rest-api/internal/pkg/model/buildable"
  "github.com/sweettea-io/rest-api/internal/pkg/service/deploysvc"
  "github.com/sweettea-io/work"
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

  // Create K8S train deploy and prep args.
  apiDeploy := kdeploy.Api{}
  apiDeployArgs := map[string]interface{}{"deployID": deployID}

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

  // Update Deploy stage to Deployed.
  if err := deploysvc.UpdateStageByID(deployID, buildable.Deployed); err != nil {
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

  return nil
}