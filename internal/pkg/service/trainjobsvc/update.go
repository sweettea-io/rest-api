package trainjobsvc

import (
  "github.com/sweettea-io/rest-api/internal/pkg/model"
  "github.com/sweettea-io/rest-api/internal/app"
  "fmt"
)

func UpdateStage(tj *model.TrainJob, stage string) error {
  if err := app.DB.Model(tj).Updates(map[string]interface{}{"stage": stage}).Error; err != nil {
    return fmt.Errorf("error updating TrainJob stage: %s", err.Error())
  }

  return nil
}

func UpdateStageByID(id uint, stage string) error {
  trainJob := model.TrainJob{}
  trainJob.ID = id
  return UpdateStage(&trainJob, stage)
}

func Fail(tj *model.TrainJob) error {
  if err := app.DB.Model(tj).Updates(map[string]interface{}{"failed": true}).Error; err != nil {
    return fmt.Errorf("error failing TrainJob: %s", err.Error())
  }

  return nil
}

func FailByID(id uint) error {
  trainJob := model.TrainJob{}
  trainJob.ID = id
  return Fail(&trainJob)
}