package clustersvc

import (
  "fmt"
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/pkg/model"
)

func Update(cluster *model.Cluster, updates *map[string]interface{}) error {
  if err := app.DB.Model(&cluster).Updates(*updates).Error; err != nil {
    return fmt.Errorf("error updating Cluster: %s", err.Error())
  }

  return nil
}