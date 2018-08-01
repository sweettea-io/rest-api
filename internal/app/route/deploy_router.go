package route

import (
  "net/http"
  "github.com/sweettea-io/rest-api/internal/app/middleware"
  "github.com/sweettea-io/rest-api/internal/app/payload"
  "github.com/sweettea-io/rest-api/internal/app/respond"
  "github.com/sweettea-io/rest-api/internal/app/errmsg"
  "github.com/sweettea-io/rest-api/internal/pkg/service/apiclustersvc"
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/pkg/service/projectsvc"
  "strings"
  "github.com/sweettea-io/rest-api/internal/pkg/util/unique"
  "github.com/sweettea-io/work"
  "github.com/sweettea-io/rest-api/internal/app/worker/jobs"
)

// ----------- ROUTER SETUP ------------

// Prefix for all routes in this file.
const DeployRoute = "/deploy"

func InitDeployRouter() {
  // Create deploy router.
  deployRouter := Router.PathPrefix(DeployRoute).Subrouter()

  // Attach Session-based auth middleware to all request handlers on this router.
  deployRouter.Use(middleware.SessionAuth)

  // Attach route handlers.
  deployRouter.HandleFunc("", CreateModelVersionHandler).Methods("POST")
}

// ----------- ROUTE HANDLERS -----------

/*
  Create a Deploy

  Method:  POST
  Route:   /deploy
  Payload:
    apiCluster string (required)
    projectNsp string (required)
    model      string (optional)
    sha        string (optional)
    envs       string (optional)
*/
func CreateDeployHandler(w http.ResponseWriter, req *http.Request) {
  // Parse & validate payload.
  var pl payload.CreateDeployPayload

  if !pl.Validate(req) {
    respond.Error(w, errmsg.InvalidPayload())
    return
  }

  // Find ApiCluster by slug.
  apiCluster, err := apiclustersvc.FromSlug(pl.ApiClusterSlug())

  if err != nil {
    app.Log.Errorln(err.Error())
    respond.Error(w, errmsg.ApiClusterNotFound())
  }

  // TODO: Query this backwards --> ModelVersion --> Model --> Project

  // Find project by namespace.
  project, err := projectsvc.FromNsp(strings.ToLower(pl.ProjectNsp), "Models", "Models.ModelVersions")

  if err != nil {
    app.Log.Errorln(err.Error())
    respond.Error(w, errmsg.ProjectNotFound())
    return
  }

  // Split provided model string into its components (<slug>:<version>)
  modelSlug, modelVersionVersion := pl.ModelBreakdown()

  // Validate model exists.
  modelExists := false
  for _, model := range project.Models {
    if model.Slug == modelSlug {
      modelExists = true
      break
    }
  }

  if !modelExists {
    app.Log.Errorf("Model(slug=%s, ProjectID=%v) not found.\n", modelSlug, project.ID)
    respond.Error(w, errmsg.ModelNotFound())
    return
  }

  // Create Uid for Deploy manually so that its available as the log stream key.
  deployUid := unique.NewUid()

  // Create args for CreateDeploy job.
  jobArgs := work.Q{
    "deployUid": deployUid,
    "apiClusterID": apiCluster,
    "modelVersionID": "TODO",
    "sha": pl.Sha,
    "envs": pl.Envs,
  }

  // Enqueue CreateDeploy job.
  if _, err := app.JobQueue.Enqueue(jobs.Names.CreateDeploy, jobArgs); err != nil {
    app.Log.Errorf("error scheduling CreateDeploy job: %s", err.Error())
    respond.Error(w, errmsg.CreateDeploySchedulingFailed())
    return
  }


  // TODO: stream training logs as response.
}