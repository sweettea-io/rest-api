package route

import (
  "net/http"
  "github.com/lib/pq"
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/app/errmsg"
  "github.com/sweettea-io/rest-api/internal/app/payload"
  "github.com/sweettea-io/rest-api/internal/app/respond"
  "github.com/sweettea-io/rest-api/internal/app/successmsg"
  "github.com/sweettea-io/rest-api/internal/pkg/service/clustersvc"
  "github.com/sweettea-io/rest-api/internal/pkg/service/usersvc"
  "github.com/sweettea-io/rest-api/internal/pkg/util/crypt"
  "github.com/sweettea-io/rest-api/internal/pkg/util/enc"
)

// ----------- ROUTER SETUP ------------

// Prefix for all routes in this file
const ClusterRoute = "/cluster"

func InitClusterRouter() {
  // Create Cluster router.
  ClusterRouter := Router.PathPrefix(ClusterRoute).Subrouter()

  // Attach route handlers.
  ClusterRouter.HandleFunc("", CreateClusterHandler).Methods("POST")
  ClusterRouter.HandleFunc("", GetClustersHandler).Methods("GET")
  ClusterRouter.HandleFunc("", UpdateClusterHandler).Methods("PUT")
  ClusterRouter.HandleFunc("", DeleteClusterHandler).Methods("DELETE")
}

// ----------- ROUTE HANDLERS -----------

/*
  Create a new Cluster (INTERNAL)

  Method:  POST
  Route:   /cluster
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
    app.Log.Errorln("error creating Cluster: invalid executor user password")
    respond.Error(w, errmsg.Unauthorized())
    return
  }

  // Only admin users can create clusters.
  if !executorUser.Admin {
    app.Log.Errorln("error creating Cluster: executor user must be an admin")
    respond.Error(w, errmsg.Unauthorized())
    return
  }

  // Create new cluster.
  cluster, err := clustersvc.Create(pl.Name, pl.Cloud, pl.State)

  if err != nil {
    app.Log.Errorln(err.Error())

    if err.(*pq.Error).Code.Name() == "unique_violation" {
      respond.Error(w, errmsg.ClusterAlreadyExists())
    } else {
      respond.Error(w, errmsg.ClusterCreationFailed())
    }

    return
  }

  respond.Created(w, enc.JSON{"cluster": cluster.AsJSON()})
}

/*
  Get Clusters by query criteria

  Method:  GET
  Route:   /cluster
*/
// TODO: Add support for query params to narrow down returned clusters.
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

  respond.Ok(w, enc.JSON{"clusters": fmtClusters})
}

/*
  Update a Cluster (INTERNAL)

  Method:  PUT
  Route:   /cluster
  Payload:
    executor_email    string (required)
    executor_password string (required)
    slug              string (required)
    name              string (required)
    cloud             string (required)
    state             string (required on all environments except 'local')
*/
func UpdateClusterHandler(w http.ResponseWriter, req *http.Request) {
  // Validate internal token.
  if internalToken := req.Header.Get(app.Config.AuthHeaderName); internalToken != app.Config.RestApiToken {
    respond.Error(w, errmsg.Unauthorized())
    return
  }

  // Parse & validate payload.
  var pl payload.UpdateClusterPayload

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
    app.Log.Errorln("error updating Cluster: invalid executor user password")
    respond.Error(w, errmsg.Unauthorized())
    return
  }

  // Only admin users can update clusters.
  if !executorUser.Admin {
    app.Log.Errorln("error updating Cluster: executor user must be an admin")
    respond.Error(w, errmsg.Unauthorized())
    return
  }

  // Find Cluster by slug.
  cluster, err := clustersvc.FromSlug(pl.Slug)

  if err != nil {
    app.Log.Errorln(err.Error())
    respond.Error(w, errmsg.ClusterNotFound())
  }

  // Update the cluster.
  if err := clustersvc.Update(cluster, pl.GetUpdates()); err != nil {
    app.Log.Errorln(err.Error())
    respond.Error(w, errmsg.ClusterUpdateFailed())
    return
  }

  respond.Ok(w, enc.JSON{"cluster": cluster.AsJSON()})
}

/*
  Delete a Cluster (INTERNAL)

  Method:  DELETE
  Route:   /cluster
  Payload:
    executor_email    string (required)
    executor_password string (required)
    slug              string (required)
*/
func DeleteClusterHandler(w http.ResponseWriter, req *http.Request) {
  // Validate internal token.
  if internalToken := req.Header.Get(app.Config.AuthHeaderName); internalToken != app.Config.RestApiToken {
    respond.Error(w, errmsg.Unauthorized())
    return
  }

  // Parse & validate payload.
  var pl payload.DeleteClusterPayload

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
    app.Log.Errorln("error deleting Cluster: invalid executor user password")
    respond.Error(w, errmsg.Unauthorized())
    return
  }

  // Only admin users can delete clusters.
  if !executorUser.Admin {
    app.Log.Errorln("error deleting Cluster: executor user must be an admin")
    respond.Error(w, errmsg.Unauthorized())
    return
  }

  // Find Cluster by slug.
  cluster, err := clustersvc.FromSlug(pl.Slug)

  if err != nil {
    app.Log.Errorln(err.Error())
    respond.Error(w, errmsg.ClusterNotFound())
  }

  // Delete the Cluster.
  if err := clustersvc.Delete(cluster); err != nil {
    app.Log.Errorln(err.Error())
    respond.Error(w, errmsg.ClusterDeletionFailed())
    return
  }

  respond.Ok(w, successmsg.ClusterDeletionSuccess)
}