package route

import (
  "net/http"
  "github.com/lib/pq"
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/app/errmsg"
  "github.com/sweettea-io/rest-api/internal/app/payload"
  "github.com/sweettea-io/rest-api/internal/app/respond"
  "github.com/sweettea-io/rest-api/internal/app/successmsg"
  "github.com/sweettea-io/rest-api/internal/pkg/service/usersvc"
  "github.com/sweettea-io/rest-api/internal/pkg/util/cloud"
  "github.com/sweettea-io/rest-api/internal/pkg/util/env"
  "github.com/sweettea-io/rest-api/internal/pkg/util/crypt"
  "github.com/sweettea-io/rest-api/internal/pkg/service/clustersvc"
  "github.com/sweettea-io/rest-api/internal/pkg/util/enc"
)

// ----------- ROUTER SETUP ------------

// Prefix for all routes in this file
const ClusterRoute = "/clusters"

func InitClusterRouter() {
  // Create Cluster router.
  ClusterRouter := Router.PathPrefix(ClusterRoute).Subrouter()

  // Attach route handlers.
  ClusterRouter.HandleFunc("", CreateClusterHandler).Methods("POST")
  ClusterRouter.HandleFunc("", GetClustersHandler).Methods("GET")
}

// ----------- ROUTE HANDLERS -----------

/*
  Create a new Cluster (INTERNAL)

  Method:  POST
  Route:   /clusters
  Payload:
    executor_email    string (required)
    executor_password string (required)
    name              string (required)
    cloud             string (required)
    state             string (required on all environments except 'local')
*/
func CreateClusterHandler(w http.ResponseWriter, req *http.Request) {
  // Validate internal token.
  if internalToken := req.Header.Get(app.Config.AuthHeaderName); internalToken != app.Config.RestApiToken {
    respond.Error(w, errmsg.Unauthorized())
    return
  }

  // Parse & validate payload.
  var pl payload.CreateClusterPayload

  // If payload is invalid, or state is empty on a non-local env, or cloud is invalid, respond with error.
  if !pl.Validate(req) || (pl.State == "" && app.Config.Env != env.Local) || !cloud.IsValidCloud(pl.Cloud) {
    respond.Error(w, errmsg.InvalidPayload())
  }

  // Get executor user by email.
  executorUser, err := usersvc.FromEmail(pl.ExecutorEmail)

  if err != nil {
    app.Log.Error(err.Error())
    respond.Error(w, errmsg.UserNotFound())
    return
  }

  // Ensure executor user's password is correct.
  if !crypt.VerifyBcrypt(pl.ExecutorPassword, executorUser.HashedPw) {
    app.Log.Errorln("error creating new Cluster: invalid executor User password")
    respond.Error(w, errmsg.Unauthorized())
    return
  }

  // Only admin users can create companies.
  if !executorUser.Admin {
    app.Log.Errorln("error creating new Cluster: executor User must be an admin an admin")
    respond.Error(w, errmsg.Unauthorized())
    return
  }

  // Create new cluster.
  cluster, err := clustersvc.Create(pl.Name, pl.Cloud, pl.State)

  if err != nil {
    app.Log.Errorf(err.Error())

    if err.(*pq.Error).Code.Name() == "unique_violation" {
      respond.Error(w, errmsg.ClusterAlreadyExists())
    } else {
      respond.Error(w, errmsg.ClusterCreationFailed())
    }

    return
  }

  // Create response payload and respond.
  respData := successmsg.ClusterCreationSuccess
  respData["slug"] = cluster.Slug

  respond.Created(w, respData)
}

/*
  Get Clusters by query criteria

  Method:  GET
  Route:   /clusters
*/
func GetClustersHandler(w http.ResponseWriter, req *http.Request) {
  // Auth request from Session token.
  _, err := usersvc.FromRequest(req)

  if err != nil {
    respond.Error(w, errmsg.Unauthorized())
    return
  }

  // Fetch all Cluster records.
  clusters := clustersvc.All()

  // Format clusters for response payload.
  var fmtClusters []enc.JSON

  for _, cluster := range clusters {
    fmtClusters = append(fmtClusters, cluster.AsJSON())
  }

  // Create response payload.
  respData := enc.JSON{
    "ok": true,
    "clusters": fmtClusters,
  }

  respond.Ok(w, respData)
}
