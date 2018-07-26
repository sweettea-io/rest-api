package clustersvc

import (
  "fmt"
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/pkg/model"
)

func Delete(cluster *model.Cluster) error {
  if err := app.DB.Delete(cluster).Error; err != nil {
    return fmt.Errorf("error deleting Cluster: %s", err.Error())
  }

  // TODO: Figure out best way to soft delete relationships as well.

  return nil
}