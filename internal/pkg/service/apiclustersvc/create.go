package apiclustersvc

import (
  "fmt"
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/pkg/model"
  "github.com/sweettea-io/rest-api/internal/pkg/util/cloud"
)

func Create(name string, cloudName string, state string) (*model.ApiCluster, error) {
  // Ensure cloud is valid.
  if !cloud.IsValidCloud(cloudName) {
    return nil, fmt.Errorf("error creating ApiCluster: \"%s\" is not a valid cloud\n", cloudName)
  }

  // Create ApiCluster model.
  apiCluster := model.ApiCluster{
    Name: name,
    Cloud: cloudName,
    State: state,
  }

  // Create record.
  if err := app.DB.Create(&apiCluster).Error; err != nil {
    return nil, fmt.Errorf("error creating ApiCluster: %s", err.Error())
  }

  return &apiCluster, nil
}

