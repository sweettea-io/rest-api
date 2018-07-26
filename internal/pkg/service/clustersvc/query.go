package clustersvc

import (
  "github.com/sweettea-io/rest-api/internal/pkg/model"
  "fmt"
  "github.com/sweettea-io/rest-api/internal/app"
)

func All() []model.Cluster {
  clusters := []model.Cluster{}
  app.DB.Find(&clusters)
  return clusters
}

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