package jobs

import (
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/pkg/service/apiclustersvc"
  "github.com/sweettea-io/rest-api/internal/pkg/service/commitsvc"
  "github.com/sweettea-io/rest-api/internal/pkg/service/deploysvc"
  "github.com/sweettea-io/rest-api/internal/pkg/service/envvarsvc"
  "github.com/sweettea-io/rest-api/internal/pkg/service/modelversionsvc"
  "github.com/sweettea-io/rest-api/internal/pkg/util/cluster"
  "github.com/sweettea-io/work"
  "github.com/sweettea-io/rest-api/internal/pkg/util/enc"
  "github.com/sweettea-io/rest-api/internal/pkg/model"
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
    logKey         (string) log key for buildable
*/
func (c *Context) CreateDeploy(job *work.Job) error {
  // Extract args from job.
  deployUid := job.ArgString("deployUid")
  deployName := job.ArgString("name")
  apiClusterID := uint(job.ArgInt64("apiClusterID"))
  modelVersionID := uint(job.ArgInt64("modelVersionID"))
  sha := job.ArgString("sha")
  envs := job.ArgString("envs")
  logKey := job.ArgString("logKey")

  if err := job.ArgError(); err != nil {
    if logKey != "" {
      return logBuildableErr(err, logKey, "Arg error occurred inside create deploy job.")
    }

    app.Log.Errorln(err.Error())
    return err
  }

  // Find ApiCluster by ID.
  apiCluster, err := apiclustersvc.FromID(apiClusterID)
  if err != nil {
    return logBuildableErr(err, logKey, "API cluster not found.")
  }

  // Find ModelVersion by ID.
  modelVersion, err := modelversionsvc.FromID(modelVersionID)
  if err != nil {
    return logBuildableErr(err, logKey, "Model version not found.")
  }

  // Store ref to project.
  project := &modelVersion.Model.Project

  // If sha provided, find Commit by that value. Otherwise, fetch latest commit from repo.
  commit, err := commitsvc.FromShaOrLatest(sha, project)
  if err != nil {
    return logBuildableErr(err, logKey, "Error finding commit sha to deploy.")
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
    return logBuildableErr(err, logKey, "Failed to upsert deploy.")
  }

  // If Deploy already exists, return an "Everything up-to-date." message.
  if !isNew {
    // TODO: stream back a success message with "Everything up-to-date."
    return nil
  }

  // Convert stringified envs into map[string]string representation.
  envsMap, err := envvarsvc.MapFromBytes([]byte(envs))
  if err != nil {
    return failDeploy(deploy.ID, err, logKey, "Failed to parse deploy environment variables.")
  }

  // Create EnvVars for this Deploy.
  if err := envvarsvc.CreateFromMap(deploy.ID, envsMap); err != nil {
    return failDeploy(deploy.ID, err, logKey, "Failed to create deploy environment variables.")
  }

  // Define args for the BuildDeploy job.
  jobArgs := work.Q{
    "resourceID": deploy.ID,
    "buildTargetSha": commit.Sha,
    "projectID": project.ID,
    "targetCluster": cluster.Api,
    "logKey": logKey,
    "followOnJob": Names.ApiDeploy,
    "followOnArgs": enc.JSON{
      "deployID": deploy.ID,
      "logKey": logKey,
    },
  }

  // Enqueue new job to build this Project for the ApiCluster.
  if _, err := app.JobQueue.Enqueue(Names.BuildDeploy, jobArgs); err != nil {
    return failDeploy(deploy.ID, err, logKey, "Failed to schedule build deploy job.")
  }

  // Update deploy stage to BuildScheduled.
  if err := deploysvc.UpdateStage(deploy, model.BuildStages.BuildScheduled); err != nil {
    return failDeploy(deploy.ID, err, logKey, "Failed to update stage of deploy.")
  }

  return nil
}