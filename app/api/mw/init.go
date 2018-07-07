package mw

import (
  "github.com/jinzhu/gorm"
  "github.com/Sirupsen/logrus"
)

// Create global vars for our db and logger so that all
// middleware functions can access these values.
var db *gorm.DB
var logger *logrus.Logger

func Init(database *gorm.DB, l *logrus.Logger) {
  db = database
  logger = l
}