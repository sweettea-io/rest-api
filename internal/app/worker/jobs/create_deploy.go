package jobs

import (
  "fmt"
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/pkg/model"
  "github.com/sweettea-io/rest-api/internal/pkg/model/buildable"
  "github.com/sweettea-io/rest-api/internal/pkg/service/apiclustersvc"
  "github.com/sweettea-io/rest-api/internal/pkg/service/commitsvc"
  "github.com/sweettea-io/rest-api/internal/pkg/service/deploysvc"
  "github.com/sweettea-io/rest-api/internal/pkg/service/envvarsvc"
  "github.com/sweettea-io/rest-api/internal/pkg/service/modelversionsvc"
  "github.com/sweettea-io/rest-api/internal/pkg/util/cluster"
  "github.com/sweettea-io/work"
  "github.com/sweettea-io/rest-api/internal/pkg/util/enc"
)

/*
  CreateDeploy handles all of the required model creation/upsertion
  leading up to a deploy to an API cluster.

  Args:
    deployUid      (string) Uid to assign to Deploy during creation
    name           (string) Name of Deploy
    apiClusterID   (uint)   ID of ApiCluster to deploy to
    modelVersionID (uint)   ID of ModelVersion to use with this deploy
    sha            (string) commit sha to deploy
    envs           (string) custom env vars to assign to this Deploy
*/
func (c *Context) CreateDeploy(job *work.Job) error {
  // Extract args from job.
  deployUid := job.ArgString("deployUid")
  deployName := job.ArgString("name")
  apiClusterID := uint(job.ArgInt64("apiClusterID"))
  modelVersionID := uint(job.ArgInt64("modelVersionID"))
  sha := job.ArgString("sha")
  envs := job.ArgString("envs")

  if err := job.ArgError(); err != nil {
    return logAndFail(err)
  }

  // Find ApiCluster by ID.
  apiCluster, err := apiclustersvc.FromID(apiClusterID)

  if err != nil {
    return logAndFail(err)
  }

  // Find ModelVersion by ID.
  modelVersion, err := modelversionsvc.FromID(modelVersionID)

  if err != nil {
    return logAndFail(err)
  }

  project := modelVersion.Model.Project
  var commit *model.Commit

  // If sha provided, find Commit by that value.
  // Otherwise, fetch the latest commit from the project's repo host.
  if sha != "" {
    var err error
    commit, err = commitsvc.FromSha(sha)

    if err != nil {
      return logAndFail(err)
    }
  } else {
    // Get host for this project.
    host := project.GetHost()
    host.Configure()

    // Get latest commit sha for project.
    latestSha, err := host.LatestSha(project.Owner(), project.Repo())

    if err != nil {
      return logAndFail(err)
    }

    // Upsert Commit for fetched sha.
    var commitUpsertErr error
    commit, err = commitsvc.Upsert(project.ID, latestSha)

    if commitUpsertErr != nil {
      return logAndFail(err)
    }
  }

  // Upsert Deploy.
  deploy, isNew, err := deploysvc.Upsert(
    commit.ID,
    modelVersion.ID,
    apiCluster.ID,
    deployUid,
    deployName,
  )

  if err != nil {
    return logAndFail(err)
  }

  // If Deploy already exists,
  if !isNew {
    // Stream back a success message with "Everything up-to-date."
    return nil
  }

  // Convert stringified envs into map[string]string representation.
  envsMap, err := envvarsvc.MapFromBytes([]byte(envs))

  if err != nil {
    return failDeploy(deploy.ID, err)
  }

  // Create EnvVars for this Deploy.
  if err := envvarsvc.CreateFromMap(deploy.ID, envsMap); err != nil {
    return failDeploy(deploy.ID, err)
  }

  // Define args for the BuildDeploy job.
  jobArgs := work.Q{
    "resourceID": deploy.ID,
    "buildTargetSha": commit.Sha,
    "projectID": project.ID,
    "targetCluster": cluster.Api,
    "followOnJob": Names.ApiDeploy,
    "followOnArgs": enc.JSON{
      "deployID": deploy.ID,
    },
  }

  // Enqueue new job to build this Project for the ApiCluster.
  if _, err := app.JobQueue.Enqueue(Names.BuildDeploy, jobArgs); err != nil {
    return failDeploy(deploy.ID, fmt.Errorf("error scheduling BuildDeploy job: %s", err.Error()))
  }

  // Update deploy stage to BuildScheduled.
  if err := deploysvc.UpdateStage(deploy, buildable.BuildScheduled); err != nil {
    return failDeploy(deploy.ID, err)
  }

  return nil
}