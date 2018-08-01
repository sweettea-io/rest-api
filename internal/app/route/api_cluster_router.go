package route

import (
  "net/http"
  "github.com/lib/pq"
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/app/errmsg"
  "github.com/sweettea-io/rest-api/internal/app/payload"
  "github.com/sweettea-io/rest-api/internal/app/respond"
  "github.com/sweettea-io/rest-api/internal/app/successmsg"
  "github.com/sweettea-io/rest-api/internal/pkg/service/apiclustersvc"
  "github.com/sweettea-io/rest-api/internal/pkg/service/usersvc"
  "github.com/sweettea-io/rest-api/internal/pkg/util/crypt"
  "github.com/sweettea-io/rest-api/internal/pkg/util/enc"
)

// ----------- ROUTER SETUP ------------

// Prefix for all routes in this file
const ApiClusterRoute = "/api_cluster"

func InitApiClusterRouter() {
  // Create apiCluster router.
  apiClusterRouter := Router.PathPrefix(ApiClusterRoute).Subrouter()

  // Attach route handlers.
  apiClusterRouter.HandleFunc("", CreateApiClusterHandler).Methods("POST")
  apiClusterRouter.HandleFunc("", GetApiClustersHandler).Methods("GET")
  apiClusterRouter.HandleFunc("", UpdateApiClusterHandler).Methods("PUT")
  apiClusterRouter.HandleFunc("", DeleteApiClusterHandler).Methods("DELETE")
}

// ----------- ROUTE HANDLERS -----------

/*
  Create a new ApiCluster (INTERNAL)

  Method:  POST
  Route:   /api_cluster
  Payload:
    executorEmail     string (required)
    executorPassword  string (required)
    name              string (required)
    cloud             string (required)
    state             string (required on all environments except 'local')
*/
func CreateApiClusterHandler(w http.ResponseWriter, req *http.Request) {
  // Validate internal token.
  if internalToken := req.Header.Get(app.Config.AuthHeaderName); internalToken != app.Config.RestApiToken {
    respond.Error(w, errmsg.Unauthorized())
    return
  }

  // Parse & validate payload.
  var pl payload.CreateApiClusterPayload

  if !pl.Validate(req) {
    respond.Error(w, errmsg.InvalidPayload())
    return
  }

  // Get executor user by email.
  executorUser, err := usersvc.FromEmail(pl.ExecutorEmail)

  if err != nil {
    app.Log.Errorln(err.Error())
    respond.Error(w, errmsg.UserNotFound())
    return
  }

  // Ensure executor user's password is correct.
  if !crypt.VerifyBcrypt(pl.ExecutorPassword, executorUser.HashedPw) {
    app.Log.Errorln("error creating ApiCluster: invalid executor user password")
    respond.Error(w, errmsg.Unauthorized())
    return
  }

  // Only admin users can create api clusters.
  if !executorUser.Admin {
    app.Log.Errorln("error creating ApiCluster: executor user must be an admin")
    respond.Error(w, errmsg.Unauthorized())
    return
  }

  // Create new ApiCluster.
  apiCluster, err := apiclustersvc.Create(pl.Name, pl.Cloud, pl.State)

  if err != nil {
    app.Log.Errorln(err.Error())
    pqError, ok := err.(*pq.Error)

    if ok && pqError.Code.Name() == "unique_violation" {
      respond.Error(w, errmsg.ApiClusterAlreadyExists())
    } else {
      respond.Error(w, errmsg.ApiClusterCreationFailed())
    }

    return
  }

  respond.Created(w, enc.JSON{"apiCluster": apiCluster.AsJSON()})
}

/*
  Get ApiClusters by query criteria

  Method:  GET
  Route:   /api_cluster
*/
// TODO: Add support for query params to narrow down returned api clusters.
func GetApiClustersHandler(w http.ResponseWriter, req *http.Request) {
  // Auth request from Session token.
  _, err := usersvc.FromRequest(req)

  if err != nil {
    respond.Error(w, errmsg.Unauthorized())
    return
  }

  // Fetch all ApiCluster records.
  apiClusters := apiclustersvc.All()

  // Format api clusters for response payload.
  var fmtApiClusters []enc.JSON
  for _, ac := range apiClusters {
    fmtApiClusters = append(fmtApiClusters, ac.AsJSON())
  }

  respond.Ok(w, enc.JSON{"apiClusters": fmtApiClusters})
}

/*
  Update a ApiCluster (INTERNAL)

  Method:  PUT
  Route:   /api_cluster
  Payload:
    executorEmail     string (required)
    executorPassword  string (required)
    slug              string (required)
    updates:          struct (optional)
      name            string (optional)
      cloud           string (optional)
      state           string (optional)
*/
func UpdateApiClusterHandler(w http.ResponseWriter, req *http.Request) {
  // Validate internal token.
  if internalToken := req.Header.Get(app.Config.AuthHeaderName); internalToken != app.Config.RestApiToken {
    respond.Error(w, errmsg.Unauthorized())
    return
  }

  // Parse & validate payload.
  var pl payload.UpdateApiClusterPayload

  if !pl.Validate(req) {
    respond.Error(w, errmsg.InvalidPayload())
    return
  }

  // Get executor user by email.
  executorUser, err := usersvc.FromEmail(pl.ExecutorEmail)

  if err != nil {
    app.Log.Errorln(err.Error())
    respond.Error(w, errmsg.UserNotFound())
    return
  }

  // Ensure executor user's password is correct.
  if !crypt.VerifyBcrypt(pl.ExecutorPassword, executorUser.HashedPw) {
    app.Log.Errorln("error updating ApiCluster: invalid executor user password")
    respond.Error(w, errmsg.Unauthorized())
    return
  }

  // Only admin users can update api clusters.
  if !executorUser.Admin {
    app.Log.Errorln("error updating ApiCluster: executor user must be an admin")
    respond.Error(w, errmsg.Unauthorized())
    return
  }

  // Find ApiCluster from slug.
  apiCluster, err := apiclustersvc.FromSlug(pl.Slug)

  if err != nil {
    app.Log.Errorln(err.Error())
    respond.Error(w, errmsg.ApiClusterNotFound())
  }

  // Update the cluster.
  if err := apiclustersvc.Update(apiCluster, pl.GetUpdates()); err != nil {
    app.Log.Errorln(err.Error())
    respond.Error(w, errmsg.ApiClusterUpdateFailed())
    return
  }

  respond.Ok(w, enc.JSON{"apiCluster": apiCluster.AsJSON()})
}

/*
  Delete a ApiCluster (INTERNAL)

  Method:  DELETE
  Route:   /api_cluster
  Payload:
    executorEmail     string (required)
    executorPassword  string (required)
    slug              string (required)
*/
func DeleteApiClusterHandler(w http.ResponseWriter, req *http.Request) {
  // Validate internal token.
  if internalToken := req.Header.Get(app.Config.AuthHeaderName); internalToken != app.Config.RestApiToken {
    respond.Error(w, errmsg.Unauthorized())
    return
  }

  // Parse & validate payload.
  var pl payload.DeleteApiClusterPayload

  if !pl.Validate(req) {
    respond.Error(w, errmsg.InvalidPayload())
    return
  }

  // Get executor user by email.
  executorUser, err := usersvc.FromEmail(pl.ExecutorEmail)

  if err != nil {
    app.Log.Errorln(err.Error())
    respond.Error(w, errmsg.UserNotFound())
    return
  }

  // Ensure executor user's password is correct.
  if !crypt.VerifyBcrypt(pl.ExecutorPassword, executorUser.HashedPw) {
    app.Log.Errorln("error deleting ApiCluster: invalid executor user password")
    respond.Error(w, errmsg.Unauthorized())
    return
  }

  // Only admin users can delete api clusters.
  if !executorUser.Admin {
    app.Log.Errorln("error deleting ApiCluster: executor user must be an admin")
    respond.Error(w, errmsg.Unauthorized())
    return
  }

  // Find ApiCluster by slug.
  apiCluster, err := apiclustersvc.FromSlug(pl.Slug)

  if err != nil {
    app.Log.Errorln(err.Error())
    respond.Error(w, errmsg.ApiClusterNotFound())
  }

  // Delete the ApiCluster.
  if err := apiclustersvc.Delete(apiCluster); err != nil {
    app.Log.Errorln(err.Error())
    respond.Error(w, errmsg.ApiClusterDeletionFailed())
    return
  }

  respond.Ok(w, successmsg.ApiClusterDeletionSuccess)
}