package projectsvc

import (
  "fmt"
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/pkg/model"
)

func Delete(project *model.Project) error {
  if err := app.DB.Delete(project).Error; err != nil {
    return fmt.Errorf("error deleting Project: %s", err.Error())
  }

  // TODO: Figure out best way to soft delete relationships as well.

  return nil
}
