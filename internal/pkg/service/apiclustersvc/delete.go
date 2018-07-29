package apiclustersvc

import (
  "fmt"
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/pkg/model"
)

func Delete(apiCluster *model.ApiCluster) error {
  if err := app.DB.Delete(apiCluster).Error; err != nil {
    return fmt.Errorf("error deleting ApiCluster: %s", err.Error())
  }

  // TODO: Figure out best way to soft delete relationships as well.

  return nil
}