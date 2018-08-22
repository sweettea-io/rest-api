package route

import (
  "net/http"
  "strings"
  "github.com/sweettea-io/rest-api/internal/app/middleware"
  "github.com/sweettea-io/rest-api/internal/app/payload"
  "github.com/sweettea-io/rest-api/internal/app/respond"
  "github.com/sweettea-io/rest-api/internal/app/errmsg"
  "github.com/sweettea-io/rest-api/internal/pkg/service/projectsvc"
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/pkg/util/unique"
  "github.com/sweettea-io/rest-api/internal/app/worker/jobs"
  "github.com/sweettea-io/rest-api/internal/pkg/util/enc"
  "github.com/sweettea-io/rest-api/internal/pkg/service/trainjobsvc"
  "github.com/sweettea-io/work"
  "github.com/sweettea-io/rest-api/internal/pkg/service/buildablesvc"
  "fmt"
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
  trainJobRouter.HandleFunc("", GetTrainJobsHandler).Methods("GET")
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
    envs        string (optional)
*/
func CreateTrainJobHandler(w http.ResponseWriter, req *http.Request) {
  // Get response streaming resources (and check if supported).
  flusher, closeNotifyCh, ok := respond.StreamResources(w)

  if !ok {
    respond.Error(w, errmsg.StreamingNotSupported())
    return
  }

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
    "sha": pl.Sha,
    "envs": pl.Envs,
  }

  // Enqueue CreateTrainJob job.
  if _, err := app.JobQueue.Enqueue(jobs.Names.CreateTrainJob, jobArgs); err != nil {
    app.Log.Errorf("error scheduling CreateTrainJob job: %s", err.Error())
    respond.Error(w, errmsg.CreateTrainJobSchedulingFailed())
    return
  }

  // TODO: prep writer with stream headers.

  // Create buildable log streamer.
  logStreamer := buildablesvc.NewLogStreamer(trainJobUid)

  // Start listening for training logs.
  go logStreamer.Watch()

  // Respond with stream of training logs.
  for {
    select {
    // Return if client closes the connection.
    case <-closeNotifyCh:
      return
    default:
      // Parse BuildableLog messages as they come in.
      log := <-logStreamer.Channel

      // If error-level log is received, fail the buildable, close the client connection, and return.
      if log.Level == "error" {
        // TODO: fail buildable and close client connection.
        return
      }

      // If buildable has completed its lifecycle, close the client connection, and return.
      if log.Complete {
        // TODO: close client connection.
        return
      }

      // Stream the log message to the client and flush data immediately.
      fmt.Fprintln(w, log.Msg)
      flusher.Flush()
    }
  }
}

/*
  Get TrainJobs by query criteria

  Method:  GET
  Route:   /train_job
*/
func GetTrainJobsHandler(w http.ResponseWriter, req *http.Request) {
  // Fetch all TrainJob records.
  trainJobs := trainjobsvc.All()

  // Format TrainJobs for response payload.
  var fmtTrainJobs []enc.JSON

  for _, tj := range trainJobs {
    fmtTrainJobs = append(fmtTrainJobs, tj.AsJSON())
  }

  respond.Ok(w, enc.JSON{"trainJobs": fmtTrainJobs})
}