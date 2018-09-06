package jobs

import (
  "fmt"
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/pkg/service/commitsvc"
  "github.com/sweettea-io/rest-api/internal/pkg/service/modelsvc"
  "github.com/sweettea-io/rest-api/internal/pkg/service/modelversionsvc"
  "github.com/sweettea-io/rest-api/internal/pkg/service/projectsvc"
  "github.com/sweettea-io/rest-api/internal/pkg/service/trainjobsvc"
  "github.com/sweettea-io/rest-api/internal/pkg/util/cluster"
  "github.com/sweettea-io/work"
  "github.com/sweettea-io/rest-api/internal/pkg/util/enc"
  m "github.com/sweettea-io/rest-api/internal/pkg/model"
)

/*
  CreateTrainJob handles all of the required model creation/upsertion
  leading up to a deploy to the Train cluster.

  Args:
    trainJobUid (string) Uid to assign to TrainJob during creation
    projectID   (uint)   ID of project associated with this TrainJob
    modelSlug   (string) slug of Model associated with this TrainJob
    sha         (string) commit sha to train with
    envs        (string) env vars to use with this TrainJob's Train Cluster deploy (json string)
    logKey      (string) log key for buildable
*/
func (c *Context) CreateTrainJob(job *work.Job) error {
  // Extract args from job.
  trainJobUid := job.ArgString("trainJobUid")
  projectID := uint(job.ArgInt64("projectID"))
  modelSlug := job.ArgString("modelSlug")
  sha := job.ArgString("sha")
  envs := job.ArgString("envs")
  logKey := job.ArgString("logKey")

  if err := job.ArgError(); err != nil {
    if logKey != "" {
      return logBuildableErr(err, logKey, "Arg error occurred inside create train job.")
    }

    app.Log.Errorln(err.Error())
    return err
  }

  // Ensure Train Cluster even exists.
  if !app.Config.TrainClusterConfigured() {
    err := fmt.Errorf("Train Cluster not yet configured.")
    return logBuildableErr(err, logKey, err.Error())
  }

  // Get project by ID.
  project, err := projectsvc.FromID(projectID)
  if err != nil {
    return logBuildableErr(err, logKey, "SweetTea project not found.")
  }

  // If sha provided, find Commit by that value. Otherwise, fetch latest commit from repo.
  commit, err := commitsvc.FromShaOrLatest(sha, project)
  if err != nil {
    return logBuildableErr(err, logKey, "Error finding commit sha to deploy.")
  }

  // Upsert Model for provided slug.
  model, err := modelsvc.Upsert(project.ID, modelSlug)
  if err != nil {
    return logBuildableErr(err, logKey, "Upserting model \"%s\" failed.", modelSlug)
  }

  // Create new ModelVersion for this model.
  modelVersion, err := modelversionsvc.Create(model.ID)
  if err != nil {
    return logBuildableErr(err, logKey, "Creating new version of model \"%s\" failed.", modelSlug)
  }

  // Create new TrainJob.
  trainJob, err := trainjobsvc.Create(trainJobUid, commit.ID, modelVersion.ID)
  if err != nil {
    return logBuildableErr(err, logKey, "Creating train job failed.")
  }

  // Follow-on job args.
  followOnArgs := &enc.JSON{
    "trainJobID": trainJob.ID,
    "envs": envs,
    "logKey": logKey,
  }

  followOnArgsStr, err := followOnArgs.AsString()
  if err != nil {
    return logBuildableErr(err, logKey, "Arg error during build deploy job scheduling.")
  }

  // Enqueue new job to build this Project for the Train Cluster.
  jobArgs := work.Q{
    "resourceID": trainJob.ID,
    "buildTargetSha": commit.Sha,
    "projectID": projectID,
    "targetCluster": cluster.Train,
    "logKey": logKey,
    "followOnJob": Names.TrainDeploy,
    "followOnArgs": followOnArgsStr,
  }

  if _, err := app.JobQueue.Enqueue(Names.BuildDeploy, jobArgs); err != nil {
    return failTrainJob(trainJob.ID, err, logKey, "Failed to schedule build deploy job.")
  }

  // Update trainJob stage to BuildScheduled.
  if err := trainjobsvc.UpdateStage(trainJob, m.BuildStages.BuildScheduled); err != nil {
    return failTrainJob(trainJob.ID, err, logKey, "Failed to update stage of train job.")
  }
  
  return nil
}