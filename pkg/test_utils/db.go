package test_utils

import (
  "github.com/jinzhu/gorm"
  "github.com/sweettea-io/rest-api/pkg/models"
)

func ClearTables(db *gorm.DB, originalLogMode bool) {
  // Turn down log mode to hide noise during this step.
  db.LogMode(false)

  // Clear all DB tables.
  db.Unscoped().Delete(&models.User{})
  db.Unscoped().Delete(&models.Session{})
  db.Unscoped().Delete(&models.Company{})
  db.Unscoped().Delete(&models.Cluster{})
  db.Unscoped().Delete(&models.Project{})
  db.Unscoped().Delete(&models.Dataset{})
  db.Unscoped().Delete(&models.Env{})
  db.Unscoped().Delete(&models.Commit{})
  db.Unscoped().Delete(&models.Deploy{})

  // Revert DB log mode to its original value.
  db.LogMode(originalLogMode)
}