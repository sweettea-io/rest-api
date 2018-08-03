package route

import (
  "net/http"
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/app/errmsg"
  "github.com/sweettea-io/rest-api/internal/app/middleware"
  "github.com/sweettea-io/rest-api/internal/app/payload"
  "github.com/sweettea-io/rest-api/internal/app/respond"
  "github.com/sweettea-io/rest-api/internal/app/worker/jobs"
  "github.com/sweettea-io/rest-api/internal/pkg/model"
  "github.com/sweettea-io/rest-api/internal/pkg/service/apiclustersvc"
  "github.com/sweettea-io/rest-api/internal/pkg/service/modelversionsvc"
  "github.com/sweettea-io/rest-api/internal/pkg/util/unique"
  "github.com/sweettea-io/work"
  "github.com/sweettea-io/rest-api/internal/pkg/service/deploysvc"
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
  deployRouter.HandleFunc("", CreateDeployHandler).Methods("POST")
  deployRouter.HandleFunc("", UpdateDeployHandler).Methods("PUT")
}

// ----------- ROUTE HANDLERS -----------

/*
  Create a Deploy

  Method:  POST
  Route:   /deploy
  Payload:
    name       string (required)
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

  // Validate deploy name availability.
  if !deploysvc.NameAvailable(pl.Name) {
    respond.Error(w, errmsg.DeployNameUnavailable())
    return
  }

  // Find ApiCluster by slug.
  apiCluster, err := apiclustersvc.FromSlug(pl.ApiClusterSlug())

  if err != nil {
    app.Log.Errorln(err.Error())
    respond.Error(w, errmsg.ApiClusterNotFound())
    return
  }

  // Find ModelVersion --> Model --> Project.
  var modelVersion *model.ModelVersion
  var mvError error

  // Use model slug and version to find ModelVersion.
  modelSlug, version := pl.ModelBreakdown()

  if version != "" {
    // If version provided, query by that.
    modelVersion, mvError = modelversionsvc.PreloadFromVersion(version, modelSlug, pl.ProjectNsp)
  } else {
    // Otherwise, take the most recently created one.
    modelVersion, mvError = modelversionsvc.PreloadLatest(modelSlug, pl.ProjectNsp)
  }

  if mvError != nil {
    app.Log.Errorln(err.Error())
    respond.Error(w, errmsg.ModelVersionNotFound())
    return
  }

  // Create Uid for Deploy manually so that its available as the log stream key.
  deployUid := unique.NewUid()

  // Create args for CreateDeploy job.
  jobArgs := work.Q{
    "deployUid": deployUid,
    "name": pl.Name,
    "apiClusterID": apiCluster.ID,
    "modelVersionID": modelVersion.ID,
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

/*
  Update a Deploy

  Method:  PIUT
  Route:   /deploy
  Payload:
    name       string (required)
    projectNsp string (required)
    model      string (optional)
    sha        string (optional)
    envs       string (optional)
*/
func UpdateDeployHandler(w http.ResponseWriter, req *http.Request) {
  // Parse & validate payload.
  var pl payload.UpdateDeployPayload

  if !pl.Validate(req) {
    respond.Error(w, errmsg.InvalidPayload())
    return
  }

  // Find Deploy by slug.
  deploy, err := deploysvc.FromSlug(pl.Slug())

  if err != nil {
    respond.Error(w, errmsg.DeployNotFound())
    return
  }

  modelVersion := deploy.ModelVersion

  // If model was provided, validate it (and its version) exists.
  if pl.Model != "" {
    var modelVersion *model.ModelVersion
    var mvError error

    // Use model slug and version to find ModelVersion.
    modelSlug, version := pl.ModelBreakdown()

    if version != "" {
      // If version provided, query by that.
      modelVersion, mvError = modelversionsvc.PreloadFromVersion(version, modelSlug, pl.ProjectNsp)
    } else {
      // Otherwise, take the most recently created one.
      modelVersion, mvError = modelversionsvc.PreloadLatest(modelSlug, pl.ProjectNsp)
    }

    if mvError != nil {
      app.Log.Errorln(err.Error())
      respond.Error(w, errmsg.ModelVersionNotFound())
      return
    }
  }

  // Create args for UpdateDeploy job.
  jobArgs := work.Q{
    "deployID": deploy.ID,
    "modelVersionID": modelVersion.ID,
    "sha": pl.Sha,
    "envs": pl.Envs,
  }

  // Enqueue UpdateDeploy job.
  if _, err := app.JobQueue.Enqueue(jobs.Names.UpdateDeploy, jobArgs); err != nil {
    app.Log.Errorf("error scheduling UpdateDeploy job: %s", err.Error())
    respond.Error(w, errmsg.UpdateDeploySchedulingFailed())
    return
  }
}