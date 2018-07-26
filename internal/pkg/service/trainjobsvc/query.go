package trainjobsvc

import (
  "github.com/sweettea-io/rest-api/internal/pkg/model"
  "github.com/sweettea-io/rest-api/internal/app"
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