package jobs

import (
  "fmt"
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/pkg/service/buildablesvc"
  "github.com/sweettea-io/rest-api/internal/pkg/util/cluster"
  "github.com/sweettea-io/work"
  "encoding/json"
  "github.com/sweettea-io/rest-api/internal/pkg/k"
  "github.com/sweettea-io/rest-api/internal/pkg/model"
)

/*
  BuildDeploy deploys the SweetTea build server to the SweetTea Build Cluster
  with instructions for building an image from a Project and ready-ing it for
  its own deploy to either the Train Cluster or an API Cluster

  Args:
    resourceID       (uint)   ID of the buildable model (TrainJob or Deploy -- based on the targetCluster)
    buildTargetSha   (string) sha to build target repo at
    projectID        (uint)   ID of the project to build
    targetCluster    (string) which cluster the project is being built for ('train' or 'api')
    followOnJob      (string) Name of job to run immediately after the build succeeds
    followOnArgs     (string) Args to pass to followOnJob
*/
func (c *Context) BuildDeploy(job *work.Job) error {
  // Extract args from job.
  resourceID := uint(job.ArgInt64("resourceID"))
  buildTargetSha := job.ArgString("buildTargetSha")
  projectID := uint(job.ArgInt64("projectID"))
  targetCluster := job.ArgString("targetCluster")
  followOnJob := job.ArgString("followOnJob")
  followOnArgs := job.ArgString("followOnArgs")

  if err := job.ArgError(); err != nil {
    return failBuildable(resourceID, targetCluster, err)
  }

  // Validate we're building for either the Train or Build Cluster.
  if targetCluster != cluster.Train && targetCluster != cluster.Api {
    err := fmt.Errorf("build deploy error: target cluster \"%s\" unsupported", targetCluster)
    return failBuildable(resourceID, targetCluster, err)
  }

  // Create K8S build deploy object and args.
  buildDeploy := k.Build{}

  bdArgs := map[string]interface{}{
    "resourceID": resourceID,
    "buildTargetSha": buildTargetSha,
    "projectID": projectID,
    "targetCluster": targetCluster,
  }

  // Initialize build deploy.
  if err := buildDeploy.Init(bdArgs); err != nil {
    return failBuildable(resourceID, targetCluster, err)
  }

  // Create deploy resources.
  if err := buildDeploy.Configure(); err != nil {
    return failBuildable(resourceID, targetCluster, err)
  }

  // Deploy to build cluster.
  if err := buildDeploy.Perform(); err != nil {
    return failBuildable(resourceID, targetCluster, err)
  }

  // Update buildable stage to Building.
  if err := buildablesvc.UpdateStage(resourceID, model.BuildStages.Building, targetCluster); err != nil {
    return failBuildable(resourceID, targetCluster, err)
  }

  // Get channel to watch for build result.
  resultCh := buildDeploy.GetResultChannel()

  // Watch build until success/failure occurs.
  go buildDeploy.Watch()
  deployResult := <-resultCh

  // Error out if build failed.
  if !deployResult.Ok {
    return failBuildable(resourceID, targetCluster, deployResult.Error)
  }

  // Unmarshal followOnArgs for target cluster deploy.
  var targetDeployArgs map[string]interface{}
  if followOnArgs != "" {
    json.Unmarshal([]byte(followOnArgs), &targetDeployArgs)
  }

  // Schedule deploy to target cluster.
  if _, err := app.JobQueue.Enqueue(followOnJob, targetDeployArgs); err != nil {
    return failBuildable(resourceID, targetCluster, fmt.Errorf("error scheduling %s job: %s", followOnJob, err.Error()))
  }

  // Update buildable stage to DeployScheduled.
  if err := buildablesvc.UpdateStage(resourceID, model.BuildStages.DeployScheduled, targetCluster); err != nil {
    return failBuildable(resourceID, targetCluster, err)
  }

  return nil
}
