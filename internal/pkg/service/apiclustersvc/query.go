package apiclustersvc

import (
  "fmt"
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/pkg/model"
)

// All returns all ApiCluster records ordered by slug.
func All() []model.ApiCluster {
  apiClusters := []model.ApiCluster{}
  app.DB.Order("slug desc").Find(&apiClusters)
  return apiClusters
}

// FromSlug attempts to find an ApiCluster record for the given slug.
// Will return an error if no record is found.
func FromSlug(slug string) (*model.ApiCluster, error) {
  // Find ApiCluster by slug.
  var apiCluster model.ApiCluster
  result := app.DB.Where(&model.ApiCluster{Slug: slug}).Find(&apiCluster)

  // Return error if not found.
  if result.RecordNotFound() {
    return nil, fmt.Errorf("ApiCluster(slug=%s) not found.\n", slug)
  }

  return &apiCluster, nil
}

// FromID attempts to find an ApiCluster record by the provided id.
// Will return an error if no record is found.
func FromID(id uint) (*model.ApiCluster, error) {
  // Find ApiCluster by ID.
  var apiCluster model.ApiCluster
  result := app.DB.First(&apiCluster, id)

  // Return error if not found.
  if result.RecordNotFound() {
    return nil, fmt.Errorf("ApiCluster(ID=%v) not found.\n", id)
  }

  return &apiCluster, nil
}