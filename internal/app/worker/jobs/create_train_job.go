package jobs

import (
  "fmt"
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/pkg/model/buildable"
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
*/
func (c *Context) CreateTrainJob(job *work.Job) error {
  // Ensure Train Cluster exists first.
  if !app.Config.TrainClusterConfigured() {
    return logAndFail(fmt.Errorf("train cluster not configured -- leaving CreateTrainJob"))
  }

  // Extract args from job.
  trainJobUid := job.ArgString("trainJobUid")
  projectID := uint(job.ArgInt64("projectID"))
  modelSlug := job.ArgString("modelSlug")
  sha := job.ArgString("sha")
  envs := job.ArgString("envs")

  if err := job.ArgError(); err != nil {
    return logAndFail(err)
  }

  // Get project by ID.
  project, err := projectsvc.FromID(projectID)

  if err != nil {
    return logAndFail(err)
  }

  var commit *m.Commit

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

  // Upsert Model for provided slug.
  model, err := modelsvc.Upsert(project.ID, modelSlug)

  if err != nil {
    return logAndFail(err)
  }

  // Create new ModelVersion for this model.
  modelVersion, err := modelversionsvc.Create(model.ID)

  if err != nil {
    return logAndFail(err)
  }

  // Create new TrainJob.
  trainJob, err := trainjobsvc.Create(trainJobUid, commit.ID, modelVersion.ID)

  if err != nil {
    return logAndFail(err)
  }

  // Enqueue new job to build this Project for the Train Cluster.
  jobArgs := work.Q{
    "resourceID": trainJob.ID,
    "buildTargetSha": sha,
    "projectID": projectID,
    "targetCluster": cluster.Train,
    "followOnJob": Names.TrainDeploy,
    "followOnArgs": enc.JSON{
      "trainJobID": trainJob.ID,
      "envs": envs,
    }.AsString(),
  }

  if _, err := app.JobQueue.Enqueue(Names.BuildDeploy, jobArgs); err != nil {
    return failTrainJob(trainJob.ID, fmt.Errorf("error scheduling BuildDeploy job: %s", err.Error()))
  }

  // Update trainJob stage to BuildScheduled.
  if err := trainjobsvc.UpdateStage(trainJob, buildable.BuildScheduled); err != nil {
    return failTrainJob(trainJob.ID, err)
  }
  
  return nil
}