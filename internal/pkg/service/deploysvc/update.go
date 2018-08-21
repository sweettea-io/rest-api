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

func Deployed(id uint, updates map[string]interface{}) error {
  // Prep Deploy model for ID.
  deploy := model.Deploy{}
  deploy.ID = id

  updates["stage"] = buildable.Deployed

  if err := app.DB.Model(&deploy).Updates(updates).Error; err != nil {
    return fmt.Errorf("error updating Deploy: %s", err.Error())
  }

  return nil
}

func RegisterLoadBalancerHost(deploy *model.Deploy, lbHost string) error {
  if err := app.DB.Model(deploy).Updates(map[string]interface{}{"LBHost": lbHost}).Error; err != nil {
    return fmt.Errorf("error updating LBHost on Deploy(id=%v): %s", deploy.ID, err.Error())
  }

  return nil
}

func Publicize(deploy *model.Deploy) error {
  if err := app.DB.Model(deploy).Updates(map[string]interface{}{"Public": true}).Error; err != nil {
    return fmt.Errorf("error registering Deploy(id=%v) as public: %s", deploy.ID, err.Error())
  }

  return nil
}