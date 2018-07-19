package route

import "net/http"

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

func CreateProjectHandler(w http.ResponseWriter, req *http.Request) {

}

func ReadProjectHandler(w http.ResponseWriter, req *http.Request) {

}

func UpdateProjectHandler(w http.ResponseWriter, req *http.Request) {

}

func DeleteProjectHandler(w http.ResponseWriter, req *http.Request) {

}