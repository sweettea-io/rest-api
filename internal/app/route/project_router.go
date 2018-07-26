package route

import (
  "net/http"
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/app/errmsg"
  "github.com/sweettea-io/rest-api/internal/app/payload"
  "github.com/sweettea-io/rest-api/internal/app/respond"
  "github.com/sweettea-io/rest-api/internal/app/successmsg"
  "github.com/sweettea-io/rest-api/internal/pkg/service/projectsvc"
  "github.com/sweettea-io/rest-api/internal/pkg/service/usersvc"
  "github.com/sweettea-io/rest-api/internal/pkg/util/enc"
  "strings"
)

// ----------- ROUTER SETUP ------------

// Prefix for all routes in this file
const ProjectRoute = "/project"

func InitProjectRouter() {
  // Create project router.
  projectRouter := Router.PathPrefix(ProjectRoute).Subrouter()

  // Attach route handlers.
  projectRouter.HandleFunc("", UpsertProjectHandler).Methods("POST")
  projectRouter.HandleFunc("", GetProjectsHandler).Methods("GET")
  projectRouter.HandleFunc("", DeleteProjectHandler).Methods("DELETE")
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
  // Auth request from Session token.
  _, err := usersvc.FromRequest(req)

  if err != nil {
    respond.Error(w, errmsg.Unauthorized())
    return
  }

  // Parse & validate payload.
  var pl payload.CreateProjectPayload

  if !pl.Validate(req) {
    respond.Error(w, errmsg.InvalidPayload())
    return
  }

  // Upsert project by namespace.
  project, isNew, err := projectsvc.UpsertByNsp(strings.ToLower(pl.Nsp))

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


/*
  Get Projects by query criteria

  Method:  GET
  Route:   /project
*/
func GetProjectsHandler(w http.ResponseWriter, req *http.Request) {
  // Auth request from Session token.
  _, err := usersvc.FromRequest(req)

  if err != nil {
    respond.Error(w, errmsg.Unauthorized())
    return
  }

  // Fetch all Project records.
  projects := projectsvc.All()

  // Format clusters for response payload.
  var fmtProjects []enc.JSON

  for _, project := range projects {
    fmtProjects = append(fmtProjects, project.AsJSON())
  }

  respond.Ok(w, enc.JSON{"projects": fmtProjects})
}

/*
  Delete a Project

  Method:  DELETE
  Route:   /project
  Payload:
    nsp    string (required)
*/
func DeleteProjectHandler(w http.ResponseWriter, req *http.Request) {
  // Auth request from Session token.
  user, err := usersvc.FromRequest(req)

  // Executor user must be an admin to delete a project.
  if err != nil || !user.Admin {
    respond.Error(w, errmsg.Unauthorized())
    return
  }

  // Parse & validate payload.
  var pl payload.DeleteProjectPayload

  if !pl.Validate(req) {
    respond.Error(w, errmsg.InvalidPayload())
    return
  }

  // Find Project by nsp.
  project, err := projectsvc.FromNsp(strings.ToLower(pl.Nsp))

  if err != nil {
    app.Log.Errorln(err.Error())
    respond.Error(w, errmsg.ProjectNotFound())
    return
  }

  // Delete the Project.
  if err := projectsvc.Delete(project); err != nil {
    app.Log.Errorln(err.Error())
    respond.Error(w, errmsg.ProjectDeletionFailed())
    return
  }

  respond.Ok(w, successmsg.ProjectDeletionSuccess)
}