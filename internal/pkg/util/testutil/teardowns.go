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
  db.Unscoped().Delete(&model.Project{})
  db.Unscoped().Delete(&model.ProjectConfig{})
  db.Unscoped().Delete(&model.Commit{})
  db.Unscoped().Delete(&model.TrainJob{})
  db.Unscoped().Delete(&model.Model{})
  db.Unscoped().Delete(&model.ModelVersion{})
  db.Unscoped().Delete(&model.Deploy{})
  db.Unscoped().Delete(&model.ApiCluster{})
  db.Unscoped().Delete(&model.EnvVar{})

  tableNames := []string{
    "users",
    "sessions",
    "projects",
    "project_configs",
    "commits",
    "train_jobs",
    "models",
    "model_versions",
    "deploys",
    "api_clusters",
    "env_vars",
  }

  for _, tableName := range tableNames {
    // Reset all primary key sequences to 1.
    db.Exec("ALTER SEQUENCE ? RESTART WITH 1;", tableName + "_id_seq")
  }

  // Revert DB log mode to its original value.
  db.LogMode(originalLogMode)
}