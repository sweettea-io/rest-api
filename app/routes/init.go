package routes

import (
  "github.com/gorilla/mux"
  "github.com/jinzhu/gorm"
)

func InitBaseRouter(db *gorm.DB) *mux.Router {
  // Create base router
  router := mux.NewRouter()

  // Create sub-routes for each model resource
  InitUserRouter(router)

  return router
}