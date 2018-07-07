package routes

import (
  "github.com/gorilla/mux"
  "net/http"
  "github.com/sweettea-io/rest-api/app"
  "github.com/sweettea-io/rest-api/defs"
  "github.com/sweettea-io/rest-api/pkg/models"
  "github.com/sweettea-io/rest-api/app/api/e"
  "github.com/sweettea-io/rest-api/app/api/resp"
  "github.com/sweettea-io/rest-api/app/api/pl"
  "github.com/sweettea-io/rest-api/pkg/utils"
  "github.com/lib/pq"
)

// ----------- ROUTER SETUP ------------

const ClusterRoute = "/clusters" // prefix for all routes in this file

func InitClusterRouter(baseRouter *mux.Router) {
  // Create Cluster router.
  ClusterRouter := baseRouter.PathPrefix(ClusterRoute).Subrouter()

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
  internalToken := req.Header.Get(defs.AuthHeaderName)

  if internalToken != app.Config.RestApiToken {
    respError(w, e.Unauthorized())
    return
  }

  // Parse & validate payload.
  var payload pl.CreateClusterPayload

  // 'state' param can only be empty on local environments.
  if !payload.Validate(req) || (payload.State == "" && app.Config.Env != utils.Envs.Local) {
    respError(w, e.InvalidPayload())
  }

  // TODO: Validate cloud payload param is one of supported values.

  // Get executor user by email.
  var executorUser models.User
  userResult := db.Where(&models.User{Email: payload.ExecutorEmail, IsDestroyed: false}).Find(&executorUser)

  if userResult.RecordNotFound() {
    respError(w, e.UserNotFound())
    return
  }

  // Ensure executor user's password is correct.
  if !utils.VerifyPw(payload.ExecutorPassword, executorUser.HashedPw) {
    respError(w, e.Unauthorized())
    return
  }

  // Only admin users can create companies.
  if !executorUser.Admin {
    logger.Errorf("Error creating new cluster: executor user %s is not an admin.\n", executorUser.Email)
    respError(w, e.Unauthorized())
    return
  }

  // Find company to attach cluster to.
  var company models.Company
  companySlug := utils.Slugify(payload.CompanyName)
  companyResult := db.Where(&models.Company{Slug: companySlug, IsDestroyed: false}).First(&company)

  // Ensure company exists.
  if companyResult.RecordNotFound() {
    respError(w, e.CompanyNotFound())
    return
  }

  // Create Cluster.
  cluster := models.Cluster{
    Name: payload.Name,
    Cloud: payload.Cloud,
    State: payload.State,
  }

  if err := db.Create(&cluster).Error; err != nil {
    logger.Errorf("Error creating new cluster: %s\n", err.Error())

    if err.(*pq.Error).Code.Name() == "unique_violation" {
      respError(w, e.ClusterAlreadyExists())
    } else {
      respError(w, e.ClusterCreationFailed())
    }

    return
  }

  // Create response payload and respond.
  respData := resp.ClusterCreationSuccess
  respData["uid"] = cluster.Uid

  respOk(w, respData)
}