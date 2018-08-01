package deploysvc

import (
  "github.com/sweettea-io/rest-api/internal/pkg/model"
  "github.com/sweettea-io/rest-api/internal/app"
  "fmt"
)

func UpdateStageByID(id uint, stage string) error {
  // Prep Deploy model for ID.
  deploy := model.Deploy{}
  deploy.ID = id

  // Update Deploy.Stage to desired stage.
  if err := app.DB.Model(&deploy).Updates(map[string]interface{}{"stage": stage}).Error; err != nil {
    return fmt.Errorf("error updating Deploy: %s", err.Error())
  }

  return nil
}
