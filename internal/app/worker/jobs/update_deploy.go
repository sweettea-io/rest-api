package jobs

import (
  "github.com/sweettea-io/work"
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/pkg/service/deploysvc"
  "github.com/sweettea-io/rest-api/internal/pkg/service/commitsvc"
  "github.com/sweettea-io/rest-api/internal/pkg/util/cluster"
  "github.com/sweettea-io/rest-api/internal/pkg/util/enc"
)

/*
  UpdateDeploy kicks off the proper jobs in order to update a Deploy,
  whether that's by building a new image or simply updating env vars.

  Args:
    deployID       (uint)   ID of Deploy to update
    modelVersionID (uint)   ID of ModelVersion to migrate to
    sha            (string) sha of Commit to migrate to
    envs           (string) custom env vars to apply to this Deploy
    logKey         (string) log key for buildable
*/
func (c *Context) UpdateDeploy(job *work.Job) error {
  // Extract args from job.
  deployID := uint(job.ArgInt64("deployID"))
  modelVersionID := uint(job.ArgInt64("modelVersionID"))
  sha := job.ArgString("sha")
  envs := job.ArgString("envs")
  logKey := job.ArgString("logKey")

  if err := job.ArgError(); err != nil {
    if logKey != "" {
      return logBuildableErr(err, logKey, "Arg error occurred inside update deploy job.")
    }

    app.Log.Errorln(err.Error())
    return err
  }

  // Find Deploy by ID.
  deploy, err := deploysvc.FromID(deployID)
  if err != nil {
    return logBuildableErr(err, logKey, "Deploy not found.")
  }

  // Get ref to this Deploy's current Commit/Project.
  commit := &deploy.Commit
  project := &commit.Project

  // If sha was provided to migrate to, fetch & upsert its Commit.
  if sha != "" {
    commit, err = commitsvc.FetchAndUpsertFromSha(project, sha)

    if err != nil {
      return logBuildableErr(err, logKey, "Failed to upsert commit with sha %s.", sha)
    }
  }

  // Bool flags to determine what exactly is being changed.
  commitChanged := deploy.Commit.ID != commit.ID
  modelVersionChanged := deploy.ModelVersion.ID != modelVersionID
  envsChanged := envs != "" // always update when any envs provided.

  // Return that everything is up to date if there's nothing to change.
  if !commitChanged && !modelVersionChanged && !envsChanged {
    // TODO: Stream everything up to date message.
    return nil
  }

  // Prep to run the ApiUpdate job.
  jobName := Names.ApiUpdate
  jobArgs := enc.JSON{
    "deployID": deploy.ID,
    "modelVersionID": modelVersionID,
    "commitID": commit.ID,
    "envs": envs,
    "logKey": logKey,
  }

  // If the commit changed, though, schedule the BuildDeploy instead.
  if commitChanged {
    jobName = Names.BuildDeploy
    jobArgs = enc.JSON{
      "resourceID": deploy.ID,
      "buildTargetSha": commit.Sha,
      "projectID": project.ID,
      "targetCluster": cluster.Api,
      "logKey": logKey,
      "followOnJob": Names.ApiUpdate,
      "followOnArgs": jobArgs.AsString(),
    }
  }

  // Schedule the appropriate job.
  if _, err := app.JobQueue.Enqueue(jobName, jobArgs); err != nil {
    return failDeploy(deploy.ID, err, logKey, "Failed to schedule %s", jobName)
  }

  return nil
}