package api

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
    calling_email    string (required)
    calling_password string (required)
    company_name     string (required)
    name             string (required)
    cloud            string (required)
    state            string (required)
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

  if !payload.Validate(req) {
    respError(w, e.InvalidPayload())
  }

  // TODO: Validate cloud payload param is one of supported values.

  // Get calling user by email.
  var callingUser models.User
  userResult := db.Where(&models.User{Email: payload.CallingEmail, IsDestroyed: false}).Find(&callingUser)

  if userResult.RecordNotFound() {
    respError(w, e.UserNotFound())
    return
  }

  // Ensure calling user's password is correct.
  if !utils.VerifyPw(payload.CallingPassword, callingUser.HashedPw) {
    respError(w, e.Unauthorized())
    return
  }

  // Only admin users can create companies.
  if !callingUser.Admin {
    logger.Errorf("Error creating new cluster: calling user %s is not an admin.\n", callingUser.Email)
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