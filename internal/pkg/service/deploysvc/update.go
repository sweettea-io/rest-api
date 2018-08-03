package deploysvc

import (
  "github.com/sweettea-io/rest-api/internal/pkg/model"
  "github.com/sweettea-io/rest-api/internal/app"
  "fmt"
  "github.com/sweettea-io/rest-api/internal/pkg/model/buildable"
)

func UpdateStage(deploy *model.Deploy, stage string) error {
  if err := app.DB.Model(deploy).Updates(map[string]interface{}{"stage": stage}).Error; err != nil {
    return fmt.Errorf("error updating Deploy stage: %s", err.Error())
  }

  return nil
}

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

func Fail(deploy *model.Deploy) error {
  if err := app.DB.Model(deploy).Updates(map[string]interface{}{"failed": true}).Error; err != nil {
    return fmt.Errorf("error failing Deploy: %s", err.Error())
  }

  return nil
}

func FailByID(id uint) error {
  deploy := model.Deploy{}
  deploy.ID = id
  return Fail(&deploy)
}

func Deployed(id uint, deploymentName string) error {
  // Prep Deploy model for ID.
  deploy := model.Deploy{}
  deploy.ID = id

  updates := map[string]interface{}{
    "stage": buildable.Deployed,
    "deployment_name": deploymentName,
  }

  if err := app.DB.Model(&deploy).Updates(updates).Error; err != nil {
    return fmt.Errorf("error updating Deploy: %s", err.Error())
  }

  return nil
}