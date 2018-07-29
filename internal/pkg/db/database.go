package db

import (
  "fmt"
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/postgres"
  "github.com/sweettea-io/rest-api/internal/pkg/model"
)


func NewConnection(url string, logMode bool) *gorm.DB {
  // Open connection to postgres.
  database, err := gorm.Open("postgres", url)

  if err != nil {
    panic(fmt.Errorf("error connecting to DB: %s", err.Error()))
  }

  // Set provided DB log mode.
  database.LogMode(logMode)

  return database
}

func Migrate(database *gorm.DB) {
  database.AutoMigrate(
    &model.User{},
    &model.Session{},
    &model.Project{},
    &model.ProjectConfig{},
    &model.Commit{},
    &model.TrainJob{},
    &model.Model{},
    &model.ModelVersion{},
    &model.Deploy{},
    &model.ApiCluster{},
    &model.EnvVar{},
  )
}