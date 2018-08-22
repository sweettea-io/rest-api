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
  streamLog, connClosedCh, ok := respond.StreamResources(w)

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
    return
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

  // Prep response headers for log streaming.
  w.Header().Set("Content-Type", "text/plain")   // only working with text logs
  w.Header().Set("Transfer-Encoding", "chunked") // we're gonna stream the logs
  w.Header().Set("X-Accel-Buffering", "no")      // prevent logs from getting backed up inside Nginx

  // Create log streamer.
  logStreamer := buildablesvc.NewLogStreamer(trainJobUid)

  // Start listening for training logs.
  // TODO: This is bad and doesn't scale...don't want an API route handler sprouting
  // goroutines in an uncapped manner like this. Create system to handle goroutine allocation
  // in a productionized and capped manner.
  go logStreamer.Watch()

  // Respond with stream of training logs.
  for {
    select {
    // Return if client closes the connection.
    case <-connClosedCh:
      return

    // Parse logs as they come in.
    default:
      log := <-logStreamer.Channel

      // Stream log message if it exists.
      if log.Msg != "" {
        streamLog(log.Msg)
      }

      // Reading logs hit unexpected error, so return.
      if log.Error != nil {
        app.Log.Errorf("Unexpected error while streaming logs: %s\n", log.Error.Error())
        return
      }

      // Log level of "error" was used somewhere during the build, so fail the TrainJob.
      if log.Failed {
        if err := trainjobsvc.FailByUid(trainJobUid); err != nil {
          app.Log.Error(err)
        }

        return
      }

      // TrainJob reached the end of its lifecycle.
      if log.Completed {
        return
      }
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