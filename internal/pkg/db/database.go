package db

import (
  "fmt"
  "github.com/jinzhu/gorm"
  "github.com/sweettea-io/rest-api/internal/pkg/model"
  _ "github.com/jinzhu/gorm/dialects/postgres"
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
    &model.Company{},
    &model.Cluster{},
    &model.Project{},
    &model.Dataset{},
    &model.Env{},
    &model.Commit{},
    &model.Deploy{},
  )
}