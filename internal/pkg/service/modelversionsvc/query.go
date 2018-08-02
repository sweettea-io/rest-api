package modelversionsvc

import (
  "fmt"
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/pkg/model"
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

func PreloadFromVersion(version string, modelSlug string, projectNsp string) (*model.ModelVersion, error) {
  // Find ModelVersion by Version.
  var modelVersion model.ModelVersion

  result := app.DB.
    Preload("Model", "slug = ?", modelSlug).
    Preload("Model.Project", "nsp = ?", projectNsp).
    Where(&model.ModelVersion{Version: version}).
    Find(&modelVersion)

  // Return error if not found.
  if result.RecordNotFound() {
    return nil, fmt.Errorf("ModelVersion(Version=%v, Model.Slug=%s) not found.\n", version, modelSlug)
  }

  return &modelVersion, nil
}

func PreloadLatest(modelSlug string, projectNsp string) (*model.ModelVersion, error) {
  // Find ModelVersion by Version.
  var modelVersion model.ModelVersion

  result := app.DB.
    Preload("Model", "slug = ?", modelSlug).
    Preload("Model.Project", "nsp = ?", projectNsp).
    Order("created_at desc").
    First(&modelVersion)

  // Return error if not found.
  if result.RecordNotFound() {
    return nil, fmt.Errorf("No ModelVersion found for Model slug %s.\n", modelSlug)
  }

  return &modelVersion, nil
}