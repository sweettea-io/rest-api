package database

import (
  "fmt"
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/postgres"
  "github.com/sweettea/rest-api/pkg/models"
)

func Connection() *gorm.DB {
  db, err := gorm.Open("postgres", "host=localhost port=5432 dbname=rest user=rest password=rest sslmode=disable")

  if err != nil {
    panic(fmt.Errorf("Error connecting to DB: %s", err))
  }

  return db
}

func Migrate() {
  db := Connection()

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