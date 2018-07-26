package clustersvc

import (
  "github.com/sweettea-io/rest-api/internal/pkg/model"
  "github.com/sweettea-io/rest-api/internal/pkg/util/cloud"
  "fmt"
  "github.com/sweettea-io/rest-api/internal/app"
)

func Create(name string, cloudName string, state string) (*model.Cluster, error) {
  // Ensure cloud is valid.
  if !cloud.IsValidCloud(cloudName) {
    return nil, fmt.Errorf("error creating Cluster: \"%s\" is not a valid cloud.\n", cloudName)
  }

  // Create Cluster model.
  cluster := model.Cluster{
    Name: name,
    Cloud: cloudName,
    State: state,
  }

  // Create record.
  if err := app.DB.Create(&cluster).Error; err != nil {
    return nil, fmt.Errorf("error creating Cluster: %s", err.Error())
  }

  return &cluster, nil
}

