package routes

import (
  "github.com/gorilla/mux"
  "github.com/jinzhu/gorm"
  "github.com/Sirupsen/logrus"
)

var db *gorm.DB
var logger *logrus.Logger

func InitBaseRouter(database *gorm.DB) *mux.Router {
  // Create base router
  router := mux.NewRouter()

  db = database
  logger = logrus.New()

  // Create sub-routes for each model resource
  InitUserRouter(router)

  return router
}