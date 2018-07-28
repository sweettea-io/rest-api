package modelsvc

import (
  "github.com/sweettea-io/rest-api/internal/pkg/model"
  "github.com/sweettea-io/rest-api/internal/app"
  "fmt"
)

func Upsert(projectID uint, slug string) (*model.Model, error) {
  // Attempt to find model by project/slug
  var m model.Model
  result := app.DB.Where(&model.Model{ProjectID: projectID, Slug: slug}).Find(&m)

  // Create new Model if one doesn't already exist.
  if result.RecordNotFound() {
    m = model.Model{
      ProjectID: projectID,
      Name: slug,
    }

    if err := app.DB.Create(&m).Error; err != nil {
      return nil, fmt.Errorf("error creating Model: %s", err.Error())
    }
  }

  return &m, nil
}
