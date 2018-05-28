package database

import (
  "fmt"
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/postgres"
  "github.com/sweettea/rest-api/pkg/models"
)

func Connection(url string) *gorm.DB {
  db, err := gorm.Open("postgres", url)

  if err != nil {
    panic(fmt.Errorf("Error connecting to DB: %s", err))
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
    &models.Team{},
    &models.Cluster{},
    &models.Bucket{},
    &models.Repo{},
    &models.Dataset{},
    &models.Env{},
    &models.Commit{},
    &models.Deployment{},
  )
}