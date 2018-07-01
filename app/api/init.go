package api

import (
  "github.com/gorilla/mux"
  "github.com/jinzhu/gorm"
  "github.com/Sirupsen/logrus"
)

// Create global vars for our db and logger so that all other
// routes in this package can reference them.
var db *gorm.DB
var logger *logrus.Logger

func CreateRouter(baseRoute string, database *gorm.DB, l *logrus.Logger) *mux.Router {
  // Assign values to our global vars declared above.
  db = database
  logger = l

  // Create base router from provided baseRoute.
  baseRouter := mux.NewRouter().PathPrefix(baseRoute).Subrouter()

  // Attach base middleware.
  baseRouter.Use(LogRequest)

  // Create route groups for each model needing a REST-ful interface.
  InitUserRouter(baseRouter)
  InitCompanyRouter(baseRouter)
  InitClusterRouter(baseRouter)

  return baseRouter
}