package clustersvc

import (
  "fmt"
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/pkg/model"
)

// All returns all Cluster records ordered by slug.
func All() []model.Cluster {
  clusters := []model.Cluster{}
  app.DB.Order("slug desc").Find(&clusters)
  return clusters
}

// FromSlug attempts to find a Cluster record for the given slug.
// Will return an error if no record is found.
func FromSlug(slug string) (*model.Cluster, error) {
  // Find Cluster by slug.
  var cluster model.Cluster
  result := app.DB.Where(&model.Cluster{Slug: slug}).Find(&cluster)

  // Return error if not found.
  if result.RecordNotFound() {
    return nil, fmt.Errorf("Cluster(slug=%s) not found.\n", slug)
  }

  return &cluster, nil
}