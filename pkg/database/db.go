package database

import (
  "fmt"
  "github.com/jinzhu/gorm"
  "github.com/sweettea-io/rest-api/pkg/models"
  _ "github.com/jinzhu/gorm/dialects/postgres"
)

func Connection(url string) *gorm.DB {
  db, err := gorm.Open("postgres", url)

  if err != nil {
    panic(fmt.Errorf("Error connecting to DB: %s", err.Error()))
  }

  return db
}

func Migrate(url string) {
  db := Connection(url)

  db.LogMode(true)

  // Auto-migrate any changes
  db.AutoMigrate(
    &models.User{},
    &models.Session{},
    &models.Company{},
    &models.Cluster{},
    &models.Project{},
    &models.Dataset{},
    &models.Env{},
    &models.Commit{},
    &models.Deploy{},
  )
}