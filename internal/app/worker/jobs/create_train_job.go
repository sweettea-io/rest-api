package jobs

import (
  "fmt"
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/pkg/service/commitsvc"
  "github.com/sweettea-io/rest-api/internal/pkg/service/modelsvc"
  "github.com/sweettea-io/rest-api/internal/pkg/service/modelversionsvc"
  "github.com/sweettea-io/rest-api/internal/pkg/service/projectsvc"
  "github.com/sweettea-io/rest-api/internal/pkg/service/trainjobsvc"
  "github.com/sweettea-io/work"
)

/*
  CreateTrainJob handles all of the required model creation/upsertion
  leading up to a deploy to the Train cluster.

  Args:
    trainJobUid (string) Uid to assign to TrainJob during creation
    projectID   (uint)   ID of project associated with this TrainJob
    modelSlug   (string) slug of Model associated with this TrainJob
*/
func (c *Context) CreateTrainJob(job *work.Job) error {
  // Ensure Train cluster exists first.
  if !app.Config.TrainClusterConfigured() {
    err := fmt.Errorf("train cluster not configured -- leaving CreateTrainJob")
    app.Log.Errorln(err.Error())
    return err
  }

  // Extract args from job.
  trainJobUid := job.ArgString("trainJobUid")
  projectID := uint(job.ArgInt64("projectID"))
  modelSlug := job.ArgString("modelSlug")

  if err := job.ArgError(); err != nil {
    app.Log.Errorln(err.Error())
    return err
  }

  // Get project by ID.
  project, err := projectsvc.FromID(projectID)

  if err != nil {
    app.Log.Errorln(err.Error())
    return err
  }

  // Get host for this project.
  host := project.GetHost(app.Config.GitHubAccessToken)

  // Get latest commit sha for project.
  sha, err := host.LatestSha(project.Owner(), project.Repo())

  if err != nil {
    app.Log.Errorln(err.Error())
    return err
  }

  // Upsert Commit for fetched sha.
  commit, err := commitsvc.Upsert(project.ID, sha)

  if err != nil {
    app.Log.Errorln(err.Error())
    return err
  }

  // Upsert Model for provided slug.
  model, err := modelsvc.Upsert(project.ID, modelSlug)

  if err != nil {
    app.Log.Errorln(err.Error())
    return err
  }

  // Create new ModelVersion for this model.
  modelVersion, err := modelversionsvc.Create(model.ID)

  if err != nil {
    app.Log.Errorln(err.Error())
    return err
  }

  // Create new TrainJob.
  trainJob, err := trainjobsvc.Create(trainJobUid, commit.ID, modelVersion.ID)

  if err != nil {
    app.Log.Errorln(err.Error())
    return err
  }

  // Enqueue TrainDeploy job.
  if _, err := app.JobQueue.Enqueue(Names.TrainDeploy, work.Q{"trainJobID": trainJob.ID}); err != nil {
    err = fmt.Errorf("error scheduling TrainDeploy job: %s", err.Error())
    app.Log.Errorln(err.Error())
    return err
  }
  
  return nil
}