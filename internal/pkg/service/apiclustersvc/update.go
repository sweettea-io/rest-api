package apiclustersvc

import (
  "fmt"
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/pkg/model"
)

func Update(apiCluster *model.ApiCluster, updates *map[string]interface{}) error {
  if err := app.DB.Model(apiCluster).Updates(*updates).Error; err != nil {
    return fmt.Errorf("error updating ApiCluster: %s", err.Error())
  }

  return nil
}