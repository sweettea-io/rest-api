package jobs

import (
  "fmt"
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/pkg/kdeploy"
  "github.com/sweettea-io/rest-api/internal/pkg/model/buildable"
  "github.com/sweettea-io/rest-api/internal/pkg/service/buildablesvc"
  "github.com/sweettea-io/rest-api/internal/pkg/util/cluster"
  "github.com/sweettea-io/work"
)

/*
  BuildDeploy deploys the SweetTea build server to the SweetTea Build Cluster
  with instructions for building an image from a Project and ready-ing it for
  its own deploy to either the Train Cluster or an API Cluster

  Args:
    resourceID     (uint)   ID of the buildable model (TrainJob or Deploy -- based on the targetCluster)
    projectID      (uint)   ID of the project to build
    targetCluster  (string) which cluster the project is being built for ('train' or 'api')
    envs           (string) json string representation of the custom env vars to use with the final target cluster deploy
*/
func (c *Context) BuildDeploy(job *work.Job) error {
  // Extract args from job.
  resourceID := uint(job.ArgInt64("resourceID"))
  projectID := uint(job.ArgInt64("projectID"))
  targetCluster := job.ArgString("targetCluster")
  envs := job.ArgString("envs")

  if err := job.ArgError(); err != nil {
    buildablesvc.Fail(resourceID, targetCluster)
    app.Log.Errorln(err.Error())
    return err
  }

  // Validate we're building for either the Train or Build Cluster.
  if targetCluster != cluster.Train && targetCluster != cluster.Api {
    err := fmt.Errorf("build deploy error: target cluster \"%s\" unsupported", targetCluster)
    buildablesvc.Fail(resourceID, targetCluster)
    app.Log.Errorln(err.Error())
    return err
  }

  // Create K8S build deploy object and args.
  buildDeploy := kdeploy.Build{}

  bdArgs := map[string]interface{}{
    "resourceID": resourceID,
    "projectID": projectID,
    "targetCluster": targetCluster,
    "envs": envs,
  }

  // Initialize build deploy.
  if err := buildDeploy.Init(bdArgs); err != nil {
    buildablesvc.Fail(resourceID, targetCluster)
    app.Log.Errorln(err.Error())
    return err
  }

  // Create deploy resources.
  if err := buildDeploy.Configure(); err != nil {
    buildablesvc.Fail(resourceID, targetCluster)
    app.Log.Errorln(err.Error())
    return err
  }

  // Deploy to build cluster.
  if err := buildDeploy.Perform(); err != nil {
    buildablesvc.Fail(resourceID, targetCluster)
    app.Log.Errorln(err.Error())
    return err
  }

  // Update buildable stage to Building.
  if err := buildablesvc.UpdateStage(resourceID, buildable.Building, targetCluster); err != nil {
    buildablesvc.Fail(resourceID, targetCluster)
    app.Log.Errorln(err.Error())
    return err
  }

  // Get channel to watch for build result.
  resultCh := buildDeploy.GetResultChannel()

  // Watch build until success/failure occurs.
  go buildDeploy.Watch()
  deployResult := <-resultCh

  // Error out if build failed.
  if !deployResult.Ok {
    buildablesvc.Fail(resourceID, targetCluster)
    app.Log.Errorf(deployResult.Error.Error())
    return deployResult.Error
  }

  // Schedule deploy to target cluster.
  targetDeployJob, targetDeployArgs := buildDeploy.NextDeploy()

  if _, err := app.JobQueue.Enqueue(targetDeployJob, targetDeployArgs); err != nil {
    app.Log.Errorf("error scheduling %s job: %s", targetDeployJob, err.Error())
    return err
  }

  // Update buildable stage to DeployScheduled.
  if err := buildablesvc.UpdateStage(resourceID, buildable.DeployScheduled, targetCluster); err != nil {
    buildablesvc.Fail(resourceID, targetCluster)
    app.Log.Errorln(err.Error())
    return err
  }

  return nil
}
