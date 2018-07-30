package jobs

import (
  "fmt"
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/pkg/kdeploy"
  "github.com/sweettea-io/rest-api/internal/pkg/util/cluster"
  "github.com/sweettea-io/work"
)

func (c *Context) BuildDeploy(job *work.Job) error {
  // Extract args from job.
  resourceID := job.ArgString("resourceID")
  projectID := uint(job.ArgInt64("projectID"))
  targetCluster := job.ArgString("targetCluster")

  if err := job.ArgError(); err != nil {
    app.Log.Errorln(err.Error())
    return err
  }

  // Validate we're building for either the Train or Build Cluster.
  if targetCluster != cluster.Train && targetCluster != cluster.Api {
    err := fmt.Errorf("build deploy error: target cluster \"%s\" unsupported", targetCluster)
    app.Log.Errorln(err.Error())
    return err
  }

  // Create K8S build deploy object.
  buildDeploy := kdeploy.Build{}

  bdArgs := map[string]interface{}{
    "resourceID": resourceID,
    "projectID": projectID,
    "targetCluster": targetCluster,
  }

  // Initialize build deploy.
  if err := buildDeploy.Init(bdArgs); err != nil {
    app.Log.Errorln(err.Error())
    return err
  }

  // Create deploy resources.
  if err := buildDeploy.Configure(); err != nil {
    app.Log.Errorln(err.Error())
    return err
  }

  // Deploy to build cluster.
  if err := buildDeploy.Perform(); err != nil {
    app.Log.Errorln(err.Error())
    return err
  }

  // Get channel to watch for deploy result.
  resultCh := buildDeploy.GetResultChannel()

  // Watch deploy until success/failure occurs.
  go buildDeploy.Watch()

  deployResult := <-resultCh

  // Error out if deploy failed.
  if !deployResult.Ok {
    app.Log.Errorf(deployResult.Error.Error())
    return deployResult.Error
  }

  targetDeploy := buildDeploy.FollowOnDeploy()

  tdArgs := map[string]interface{}{
    "resourceID": resourceID,
  }

  // Initialize deploy to target cluster.
  if err := targetDeploy.Init(tdArgs); err != nil {
    app.Log.Errorln(err.Error())
    return err
  }

  // Create deploy resources.
  if err := targetDeploy.Configure(); err != nil {
    app.Log.Errorln(err.Error())
    return err
  }

  // Deploy to target cluster.
  if err := targetDeploy.Deploy(); err != nil {
    app.Log.Errorln(err.Error())
    return err
  }

  return nil
}
