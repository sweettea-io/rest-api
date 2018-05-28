package routes

import (
  "github.com/gorilla/mux"
  "github.com/jinzhu/gorm"
  "github.com/Sirupsen/logrus"
)

// Create global vars for our db and logger so that all other
// routes in this package can reference them.
var db *gorm.DB
var logger *logrus.Logger

func CreateRouter(baseRoute string, database *gorm.DB, logr *logrus.Logger) *mux.Router {
  // Create base router from provided baseRoute.
  baseRouter := mux.NewRouter().PathPrefix(baseRoute).Subrouter()

  // Assign these to our global vars declared above.
  db = database
  logger = logr

  // Create route groups for each model needing a RESTful interface.
  InitUserRouter(baseRouter)

  return baseRouter
}