package trainjobsvc

import (
  "github.com/sweettea-io/rest-api/internal/pkg/model"
  "github.com/sweettea-io/rest-api/internal/app"
  "fmt"
)

func UpdateStageByID(id uint, stage string) error {
  // Prep TrainJob model for ID.
  trainJob := model.TrainJob{}
  trainJob.ID = id

  // Update TrainJob.Stage to desired stage.
  if err := app.DB.Model(&trainJob).Updates(map[string]interface{}{"stage": stage}).Error; err != nil {
    return fmt.Errorf("error updating TrainJob: %s", err.Error())
  }

  return nil
}