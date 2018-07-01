package api

import (
  "github.com/gorilla/mux"
  "net/http"
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
    name             string (required)
    cloud            string (required)
    state            string (required)
*/
func CreateClusterHandler(w http.ResponseWriter, req *http.Request) {

}