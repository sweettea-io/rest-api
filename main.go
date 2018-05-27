package main

import (
  "github.com/sweettea/rest-api/pkg/models"
  "github.com/sweettea/rest-api/pkg/database"
)

func main() {
  db := database.Connection()

  // Create new user
  team := &models.Team{Name: "Gabs Team"}

  db.Create(&team)
}