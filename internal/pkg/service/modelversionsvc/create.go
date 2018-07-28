package modelversionsvc

import (
  "github.com/sweettea-io/rest-api/internal/pkg/model"
  "github.com/sweettea-io/rest-api/internal/app"
  "fmt"
)

func Create(modelID uint) (*model.ModelVersion, error) {
  // Create ModelVersion model.
  modelVersion := model.ModelVersion{ModelID: modelID}

  // Create record.
  if err := app.DB.Create(&modelVersion).Error; err != nil {
    return nil, fmt.Errorf("error creating ModelVersion: %s", err.Error())
  }

  return &modelVersion, nil
}