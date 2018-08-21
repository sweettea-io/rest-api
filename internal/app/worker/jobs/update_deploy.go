package jobs

import (
  "github.com/sweettea-io/work"
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/pkg/service/deploysvc"
  "github.com/sweettea-io/rest-api/internal/pkg/service/commitsvc"
  "github.com/sweettea-io/rest-api/internal/pkg/util/cluster"
  "github.com/sweettea-io/rest-api/internal/pkg/util/enc"
  "fmt"
)

/*
  UpdateDeploy kicks off the proper jobs in order to update a Deploy,
  whether that's by building a new image or simply updating env vars.

  Args:
    deployID       (uint)   ID of Deploy to update
    modelVersionID (uint)   ID of ModelVersion to migrate to
    sha            (string) sha of Commit to migrate to
    envs           (string) custom env vars to apply to this Deploy
*/
func (c *Context) UpdateDeploy(job *work.Job) error {
  // Extract args from job.
  deployID := uint(job.ArgInt64("deployID"))
  modelVersionID := uint(job.ArgInt64("modelVersionID"))
  sha := job.ArgString("sha")
  envs := job.ArgString("envs")

  if err := job.ArgError(); err != nil {
    return logAndFail(err)
  }

  // Find Deploy by ID.
  deploy, err := deploysvc.FromID(deployID)

  if err != nil {
    return logAndFail(err)
  }

  // Get ref to this Deploy's current Commit/Project.
  commit := &deploy.Commit
  project := commit.Project

  // If sha was provided to migrate to, fetch & upsert its Commit.
  if sha != "" {
    var err error
    host := project.GetHost()
    host.Configure()

    if sha == "latest" {
      sha, err = host.LatestSha(project.Owner(), project.Repo())
    } else {
      err = host.EnsureCommitExists(project.Owner(), project.Repo(), sha)
    }

    if err != nil {
      return logAndFail(err)
    }

    commit, err = commitsvc.Upsert(project.ID, sha)

    if err != nil {
      return logAndFail(err)
    }
  }

  // Bool flags to determine what exactly is being changed.
  commitChanged := deploy.Commit.ID != commit.ID
  modelVersionChanged := deploy.ModelVersion.ID != modelVersionID
  envsChanged := envs != "" // always update when any envs provided.

  // Return that everything is up to date if there's nothing to change.
  if !commitChanged && !modelVersionChanged && !envsChanged {
    return nil
  }

  // Prep to run the ApiUpdate job.
  jobName := Names.ApiUpdate
  jobArgs := enc.JSON{
    "deployID": deploy.ID,
    "modelVersionID": modelVersionID,
    "commitID": commit.ID,
    "envs": envs,
  }

  // If the commit changed, though, schedule the BuildDeploy instead.
  if commitChanged {
    jobName = Names.BuildDeploy
    jobArgs = enc.JSON{
      "resourceID": deploy.ID,
      "buildTargetSha": commit.Sha,
      "projectID": project.ID,
      "targetCluster": cluster.Api,
      "followOnJob": Names.ApiUpdate,
      "followOnArgs": jobArgs.AsString(),
    }
  }

  // Schedule the appropriate job.
  if _, err := app.JobQueue.Enqueue(jobName, jobArgs); err != nil {
    return failDeploy(deploy.ID, fmt.Errorf("error scheduling %s job: %s", jobName, err.Error()))
  }

  return nil
}