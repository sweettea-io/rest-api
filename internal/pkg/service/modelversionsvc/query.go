package modelversionsvc

import (
  "github.com/sweettea-io/rest-api/internal/pkg/model"
  "github.com/sweettea-io/rest-api/internal/app"
  "fmt"
)

// FromID attempts to find a ModelVersion record by the provided id.
// Will return an error if no record is found.
func FromID(id uint) (*model.ModelVersion, error) {
  // Find ModelVersion by ID.
  var modelVersion model.ModelVersion
  result := app.DB.
    Preload("Model").
    Preload("Model.Project").
    First(&modelVersion, id)

  // Return error if not found.
  if result.RecordNotFound() {
    return nil, fmt.Errorf("ModelVersion(ID=%v) not found.\n", id)
  }

  return &modelVersion, nil
}