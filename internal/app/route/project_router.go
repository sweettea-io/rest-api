package route

import (
  "net/http"
  "github.com/sweettea-io/rest-api/internal/app/respond"
  "github.com/sweettea-io/rest-api/internal/app/errmsg"
  "github.com/sweettea-io/rest-api/internal/app/payload"
  "github.com/sweettea-io/rest-api/internal/pkg/service/projectsvc"
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/pkg/util/enc"
  "github.com/sweettea-io/rest-api/internal/app/middleware"
)

// ----------- ROUTER SETUP ------------

// Prefix for all routes in this file
const ProjectRoute = "/project"

func InitProjectRouter() {
  // Create project router.
  projectRouter := Router.PathPrefix(ProjectRoute).Subrouter()

  // Add Session token auth requirement to all routes in this handler.
  projectRouter.Use(middleware.SessionAuth)

  // Attach route handlers.
  projectRouter.HandleFunc("", UpsertProjectHandler).Methods("POST")
}

// ----------- ROUTE HANDLERS -----------

/*
  Upsert a Project

  Method:  POST
  Route:   /project
  Payload:
    nsp    string (required)
*/
func UpsertProjectHandler(w http.ResponseWriter, req *http.Request) {
  // Parse & validate payload.
  var pl payload.CreateProjectPayload

  if !pl.Validate(req) {
    respond.Error(w, errmsg.InvalidPayload())
    return
  }

  // Upsert project by namespace.
  project, isNew, err := projectsvc.UpsertByNsp(pl.Nsp)

  if err != nil {
    app.Log.Errorln(err.Error())
    respond.Error(w, errmsg.ProjectUpsertionFailed())
    return
  }

  // Return unavailable error if Project already exists.
  if !isNew {
    respond.Error(w, errmsg.ProjectNspUnavailable())
    return
  }

  respond.Created(w, enc.JSON{"project": project.AsJSON()})
}