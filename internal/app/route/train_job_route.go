package route

import (
  "net/http"
  "github.com/sweettea-io/rest-api/internal/app/middleware"
  "github.com/sweettea-io/rest-api/internal/app/payload"
  "github.com/sweettea-io/rest-api/internal/app/respond"
  "github.com/sweettea-io/rest-api/internal/app/errmsg"
  "github.com/sweettea-io/rest-api/internal/pkg/service/projectsvc"
  "strings"
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/pkg/util/unique"
  "github.com/sweettea-io/work"
  "github.com/sweettea-io/rest-api/internal/app/worker/jobs"
)

// ----------- ROUTER SETUP ------------

// Prefix for all routes in this file
const TrainJobRoute = "/train_job"

func InitTrainJobRouter() {
  // Create trainJob router.
  trainJobRouter := Router.PathPrefix(TrainJobRoute).Subrouter()

  // Attach Session-based auth middleware to all request handlers on this router.
  trainJobRouter.Use(middleware.SessionAuth)

  // Attach route handlers.
  trainJobRouter.HandleFunc("", CreateTrainJobHandler).Methods("POST")
}

// ----------- ROUTE HANDLERS -----------

/*
  Create a TrainJob
  If successful, will enqueue a CreateTrainJob job.

  Method:  POST
  Route:   /train_job
  Payload:
    projectNsp  string (required)
    modelName   string (required)
*/

func CreateTrainJobHandler(w http.ResponseWriter, req *http.Request) {
  // Parse & validate payload.
  var pl payload.CreateTrainJobPayload

  if !pl.Validate(req) {
    respond.Error(w, errmsg.InvalidPayload())
    return
  }

  // Ensure Train cluster exists.
  if !app.Config.TrainClusterConfigured() {
    respond.Error(w, errmsg.TrainClusterNotConfigured())
  }

  // Find project by namespace.
  project, err := projectsvc.FromNsp(strings.ToLower(pl.ProjectNsp))

  if err != nil {
    app.Log.Errorln(err.Error())
    respond.Error(w, errmsg.ProjectNotFound())
    return
  }

  // Create Uid for TrainJob manually so that its available as the log stream key.
  trainJobUid := unique.NewUid()

  // Create args for CreateTrainJob job.
  jobArgs := work.Q{
    "trainJobUid": trainJobUid,
    "projectID": project.ID,
    "modelSlug": pl.ModelSlug(),
  }

  // Enqueue CreateTrainJob job.
  if _, err := app.JobQueue.Enqueue(jobs.Names.CreateTrainJob, jobArgs); err != nil {
    app.Log.Errorf("Error enqueueing CreateTrainJob: %s", err.Error())
    respond.Error(w, errmsg.CreateTrainJobSchedulingFailed())
    return
  }

  // TODO: stream training logs as response.
}