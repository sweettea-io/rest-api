package route

import (
  "net/http"
  "github.com/sweettea-io/rest-api/internal/pkg/service/usersvc"
  "github.com/sweettea-io/rest-api/internal/app/respond"
  "github.com/sweettea-io/rest-api/internal/app/errmsg"
  "github.com/sweettea-io/rest-api/internal/app/payload"
  "github.com/sweettea-io/rest-api/internal/pkg/util/str"
  "github.com/sweettea-io/rest-api/internal/pkg/service/companysvc"
  "github.com/sweettea-io/rest-api/internal/pkg/model"
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/lib/pq"
  "github.com/sweettea-io/rest-api/internal/app/successmsg"
)

// ----------- ROUTER SETUP ------------

// Prefix for all routes in this file
const ProjectRoute = "/projects"

func InitProjectRouter() {
  // Create project router.
  projectRouter := Router.PathPrefix(ProjectRoute).Subrouter()

  // Attach route handlers.
  projectRouter.HandleFunc("", CreateProjectHandler).Methods("POST")
  projectRouter.HandleFunc("", ReadProjectHandler).Methods("GET")
  projectRouter.HandleFunc("", UpdateProjectHandler).Methods("UPDATE")
  projectRouter.HandleFunc("", DeleteProjectHandler).Methods("DELETE")
}

/*
  Create a new Project

  Method:  POST
  Route:   /projects
  Payload:
    name    string (required)
    company string (required)
*/
func CreateProjectHandler(w http.ResponseWriter, req *http.Request) {
  // Find user from request session.
  _, err := usersvc.FromRequest(req)

  if err != nil {
    respond.Error(w, errmsg.UserNotFound())
  }

  // Parse & validate payload.
  var pl payload.CreateProjectPayload

  if !pl.Validate(req) {
    respond.Error(w, errmsg.InvalidPayload())
    return
  }

  // Ensure company associated with project exists.
  company, err := companysvc.FindBySlug(str.Slugify(pl.Company))

  if err != nil {
    respond.Error(w, errmsg.CompanyNotFound())
    return
  }

  // Create new Project.
  project := model.Project{
    Name: pl.Name,
    CompanyID: company.ID,
  }

  if err := app.DB.Create(&project).Error; err != nil {
    app.Log.Errorf("Error creating new project: %s\n", err.Error())

    if err.(*pq.Error).Code.Name() == "unique_violation" {
      respond.Error(w, errmsg.ProjectNotAvailable())
    } else {
      respond.Error(w, errmsg.ProjectCreationFailed())
    }

    return
  }

  // Create response pl and respond.
  respData := successmsg.ProjectCreationSuccess
  respData["uid"] = project.Uid

  respond.Created(w, respData)
}

func ReadProjectHandler(w http.ResponseWriter, req *http.Request) {
  // Find user from request session.
  _, err := usersvc.FromRequest(req)

  if err != nil {
    respond.Error(w, errmsg.UserNotFound())
  }

  // Do user shit
}

func UpdateProjectHandler(w http.ResponseWriter, req *http.Request) {
  // Find user from request session.
  user, err := usersvc.FromRequest(req)

  if err != nil {
    respond.Error(w, errmsg.UserNotFound())
  }

  // Do user shit
}

func DeleteProjectHandler(w http.ResponseWriter, req *http.Request) {
  // Find user from request session.
  user, err := usersvc.FromRequest(req)

  if err != nil {
    respond.Error(w, errmsg.UserNotFound())
  }

  // Do user shit
}