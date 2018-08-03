package jobs

import (
  "github.com/sweettea-io/work"
  "github.com/sweettea-io/rest-api/internal/pkg/service/deploysvc"
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/pkg/kdeploy"
  "github.com/sweettea-io/rest-api/internal/pkg/service/commitsvc"
  "github.com/sweettea-io/rest-api/internal/pkg/service/modelversionsvc"
  "github.com/sweettea-io/rest-api/internal/pkg/util/maputil"
  "github.com/sweettea-io/rest-api/internal/pkg/service/envvarsvc"
)

/*
  ApiUpdate updates a deployment for a Deploy and migrates its parent models on success (if needed).

  Args:
    deployID       (uint)   ID of the Deploy model
    modelVersionID (uint)   ModelVersion to migrate Deploy
    commitID       (uint)   Commit to migrate Deploy to
    envs           (string) stringified env var updates to apply
*/
func (c *Context) ApiUpdate(job *work.Job) error {
  // Extract args from job.
  deployID := uint(job.ArgInt64("deployID"))
  modelVersionID := uint(job.ArgInt64("modelVersionID"))
  commitID := uint(job.ArgInt64("commitID"))
  envs := job.ArgString("envs")

  if err := job.ArgError(); err != nil {
    deploysvc.FailByID(deployID)
    app.Log.Errorln(err.Error())
    return err
  }

  // Find Deploy by ID.
  deploy, err := deploysvc.FromID(deployID)

  if err != nil {
    deploysvc.FailByID(deployID)
    app.Log.Errorln(err.Error())
    return err
  }

  // Get current Commit & ModelVersion
  commit := &deploy.Commit
  modelVersion := &deploy.ModelVersion

  // Bool flags to determine what exactly is being changed.
  updatingCommit := commit.ID != commitID
  updatingModelVersion := modelVersion.ID != modelVersionID

  // Get new commit if being updated.
  if updatingCommit {
    var err error
    commit, err = commitsvc.FromID(commitID)

    if err != nil {
      deploysvc.FailByID(deployID)
      app.Log.Errorln(err.Error())
      return err
    }
  }

  // Get new modelVersion if being updated.
  if updatingModelVersion {
    var err error
    modelVersion, err = modelversionsvc.FromID(modelVersionID)

    if err != nil {
      deploysvc.FailByID(deployID)
      app.Log.Errorln(err.Error())
      return err
    }
  }

  // Convert stringified envs into map[string]string representation.
  envUpdates, err := envvarsvc.MapFromBytes([]byte(envs))

  if err != nil {
    deploysvc.FailByID(deployID)
    app.Log.Errorln(err.Error())
    return err
  }

  // Merge new envs on top of existing ones.
  allCustomEnvs := maputil.MergeMaps(envvarsvc.GetMap(deployID), envUpdates)

  // Create K8S API deploy and prep args.
  apiDeploy := kdeploy.Api{}
  apiDeployArgs := map[string]interface{}{
    "deploy": deploy,
    "commit": commit,
    "modelVersion": modelVersion,
    "customEnvs": allCustomEnvs,
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

  // Create map of updates to apply to Deploy now that it has succeeded.
  updates := map[string]interface{}{}

  if updatingCommit {
    updates["commit_id"] = commit.ID
  }

  if updatingModelVersion {
    updates["model_version_id"] = modelVersion.ID
  }

  // Update Deploy stage to Deployed and apply updates.
  if err := deploysvc.Deployed(deployID, updates); err != nil {
    deploysvc.FailByID(deployID)
    app.Log.Errorln(err.Error())
    return err
  }

  // Upsert all envs that could have been changed.
  if err := envvarsvc.UpsertFromMap(deployID, envUpdates); err != nil {
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
