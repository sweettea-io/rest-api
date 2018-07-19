package route

import (
  "net/http"
  "github.com/lib/pq"
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/app/errmsg"
  "github.com/sweettea-io/rest-api/internal/app/payload"
  "github.com/sweettea-io/rest-api/internal/app/respond"
  "github.com/sweettea-io/rest-api/internal/app/successmsg"
  "github.com/sweettea-io/rest-api/internal/pkg/model"
  "github.com/sweettea-io/rest-api/internal/pkg/util/cloud"
  "github.com/sweettea-io/rest-api/internal/pkg/util/env"
  "github.com/sweettea-io/rest-api/internal/pkg/util/crypt"
  "github.com/sweettea-io/rest-api/internal/pkg/util/str"
)

// ----------- ROUTER SETUP ------------

const ClusterRoute = "/clusters" // prefix for all routes in this file

func InitClusterRouter() {
  // Create Cluster router.
  ClusterRouter := Router.PathPrefix(ClusterRoute).Subrouter()

  // Attach route handlers.
  ClusterRouter.HandleFunc("", CreateClusterHandler).Methods("POST")
}

// ----------- ROUTE HANDLERS -----------

/*
  Create a new Cluster (INTERNAL)

  Method:  POST
  Route:   /clusters
  Payload:
    executor_email    string (required)
    executor_password string (required)
    company_name      string (required)
    name              string (required)
    cloud             string (required)
    state             string (required on all environments except 'local')
*/
func CreateClusterHandler(w http.ResponseWriter, req *http.Request) {
  // Validate internal token.
  internalToken := req.Header.Get(app.Config.AuthHeaderName)

  if internalToken != app.Config.RestApiToken {
    respond.Error(w, errmsg.Unauthorized())
    return
  }

  // Parse & validate payload.
  var pl payload.CreateClusterPayload

  if !pl.Validate(req) ||
    (pl.State == "" && app.Config.Env != env.Local) || // state required on all non-local envs
    (!cloud.IsValidCloud(pl.Cloud)) {

    respond.Error(w, errmsg.InvalidPayload())
  }

  // Get executor user by email.
  var executorUser model.User
  userResult := app.DB.Where(&model.User{Email: pl.ExecutorEmail}).Find(&executorUser)

  if userResult.RecordNotFound() {
    respond.Error(w, errmsg.UserNotFound())
    return
  }

  // Ensure executor user's password is correct.
  if !crypt.VerifyBcrypt(pl.ExecutorPassword, executorUser.HashedPw) {
    respond.Error(w, errmsg.Unauthorized())
    return
  }

  // Only admin users can create companies.
  if !executorUser.Admin {
    app.Log.Errorf("Error creating new cluster: executor user %s is not an admin.\n", executorUser.Email)
    respond.Error(w, errmsg.Unauthorized())
    return
  }

  // Find company to attach cluster to.
  var company model.Company
  companySlug := str.Slugify(pl.CompanyName)
  companyResult := app.DB.Where(&model.Company{Slug: companySlug}).First(&company)

  // Ensure company exists.
  if companyResult.RecordNotFound() {
    respond.Error(w, errmsg.CompanyNotFound())
    return
  }

  // Create Cluster.
  cluster := model.Cluster{
    Name: pl.Name,
    Cloud: pl.Cloud,
    State: pl.State,
  }

  if err := app.DB.Create(&cluster).Error; err != nil {
    app.Log.Errorf("Error creating new cluster: %s\n", err.Error())

    if err.(*pq.Error).Code.Name() == "unique_violation" {
      respond.Error(w, errmsg.ClusterAlreadyExists())
    } else {
      respond.Error(w, errmsg.ClusterCreationFailed())
    }

    return
  }

  // Create response pl and respond.
  respData := successmsg.ClusterCreationSuccess
  respData["uid"] = cluster.Uid

  respond.Created(w, respData)
}