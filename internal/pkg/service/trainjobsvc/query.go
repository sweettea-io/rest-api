package trainjobsvc

import (
  "github.com/sweettea-io/rest-api/internal/pkg/model"
  "github.com/sweettea-io/rest-api/internal/app"
  "fmt"
)

// All returns all TrainJob records ordered by time of creation.
func All() []model.TrainJob {
  trainJobs := []model.TrainJob{}

  // Find TrainJobs & eager load relationships.
  app.DB.
    Preload("Commit").
    Preload("ModelVersion").
    Preload("ModelVersion.Model").
    Order("created_at desc").
    Find(&trainJobs)

  return trainJobs
}

// FromID attempts to find a TrainJob record by the provided id.
// Will return an error if no record is found.
func FromID(id uint) (*model.TrainJob, error) {
  // Find TrainJob by ID.
  var trainJob model.TrainJob
  result := app.DB.
    Preload("Commit").
    Preload("ModelVersion").
    First(&trainJob, id)

  // Return error if not found.
  if result.RecordNotFound() {
    return nil, fmt.Errorf("TrainJob(ID=%v) not found.\n", id)
  }

  return &trainJob, nil
}