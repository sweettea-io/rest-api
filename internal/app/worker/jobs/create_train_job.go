package jobs

import (
  "fmt"
  "github.com/sweettea-io/work"
  "github.com/sweettea-io/rest-api/internal/pkg/service/commitsvc"
  "github.com/sweettea-io/rest-api/internal/pkg/service/modelsvc"
  "github.com/sweettea-io/rest-api/internal/pkg/service/modelversionsvc"
  "github.com/sweettea-io/rest-api/internal/pkg/service/projectsvc"
  "github.com/sweettea-io/rest-api/internal/pkg/service/trainjobsvc"
)

/*
  CreateTrainJob handles all of the required model creation/upsertion
  leading up to a deploy to the Train cluster.

  Args:
    trainJobUid (string) Uid to assign to TrainJob during creation
    projectID   (uint)   ID of project associated with this TrainJob
    modelSlug   (string) slug of Model associated with this TrainJob
*/
// TODO: Figure out if you need to re-initialize app in order to be able to use the svc's that reference app.Stuff.
func (c *Context) CreateTrainJob(job *work.Job) error {
  // Extract args from job.
  trainJobUid := job.ArgString("trainJobUid")
  projectID := uint(job.ArgInt64("modelID"))
  modelSlug := job.ArgString("modelSlug")

  if err := job.ArgError(); err != nil {
    c.Log.Errorln(err.Error())
    return err
  }

  // Get project by ID.
  project, err := projectsvc.FromID(projectID)

  if err != nil {
    c.Log.Errorln(err.Error())
    return err
  }

  // Get latest commit sha for project.
  sha, err := project.GetHost().LatestSha()

  if err != nil {
    c.Log.Errorln(err.Error())
    return err
  }

  // Upsert Commit for fetched sha.
  commit, err := commitsvc.Upsert(project.ID, sha)

  if err != nil {
    c.Log.Errorln(err.Error())
    return err
  }

  // Upsert Model for provided slug.
  model, err := modelsvc.Upsert(project.ID, modelSlug)

  if err != nil {
    c.Log.Errorln(err.Error())
    return err
  }

  // Create new ModelVersion for this model.
  modelVersion, err := modelversionsvc.Create(model.ID)

  if err != nil {
    c.Log.Errorln(err.Error())
    return err
  }

  // Create new TrainJob.
  trainJob, err := trainjobsvc.Create(trainJobUid, commit.ID, modelVersion.ID)

  if err != nil {
    c.Log.Errorln(err.Error())
    return err
  }

  // Enqueue TrainDeploy job.
  if _, err := c.JobQueue.Enqueue(Names.TrainDeploy, work.Q{"trainJobID": trainJob.ID}); err != nil {
    err = fmt.Errorf("error scheduling TrainDeploy job: %s", err.Error())
    c.Log.Errorln(err.Error())
    return err
  }
  
  return nil
}