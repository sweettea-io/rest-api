package trainjobsvc

import (
  "github.com/sweettea-io/rest-api/internal/pkg/model"
  "github.com/sweettea-io/rest-api/internal/app"
  "fmt"
  "github.com/sweettea-io/rest-api/internal/pkg/util/unique"
)

func Create(uid string, commitID uint, modelVersionID uint) (*model.TrainJob, error) {
  // Generate new uid if empty one provided.
  if uid == "" {
    uid = unique.NewUid()
  }

  // Create TrainJob model.
  trainJob := model.TrainJob{
    Uid: uid,
    CommitID: commitID,
    ModelVersionID: modelVersionID,
  }

  // Create record.
  if err := app.DB.Create(&trainJob).Error; err != nil {
    return nil, fmt.Errorf("error creating TrainJob: %s", err.Error())
  }

  return &trainJob, nil
}