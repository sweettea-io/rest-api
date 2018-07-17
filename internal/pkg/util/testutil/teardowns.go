package testutil

import (
  "github.com/jinzhu/gorm"
  "github.com/sweettea-io/rest-api/internal/pkg/model"
)

func ClearTables(db *gorm.DB, originalLogMode bool) {
  // Turn down log mode to hide noise during this step.
  db.LogMode(false)

  // Clear all DB tables.
  db.Unscoped().Delete(&model.User{})
  db.Unscoped().Delete(&model.Session{})
  db.Unscoped().Delete(&model.Company{})
  db.Unscoped().Delete(&model.Cluster{})
  db.Unscoped().Delete(&model.Project{})
  db.Unscoped().Delete(&model.Dataset{})
  db.Unscoped().Delete(&model.Env{})
  db.Unscoped().Delete(&model.Commit{})
  db.Unscoped().Delete(&model.Deploy{})

  tableNames := []string{
    "users",
    "sessions",
    "companies",
    "clusters",
    "projects",
    "datasets",
    "envs",
    "commits",
    "deploys",
  }

  for _, tableName := range tableNames {
    // Reset all primary key sequences to 1.
    db.Exec("ALTER SEQUENCE ? RESTART WITH 1;", tableName + "_id_seq")
  }

  // Revert DB log mode to its original value.
  db.LogMode(originalLogMode)
}