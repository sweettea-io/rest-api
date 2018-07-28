package commitsvc

import (
  "github.com/sweettea-io/rest-api/internal/pkg/model"
  "github.com/sweettea-io/rest-api/internal/app"
  "fmt"
)

func Upsert(projectID uint, sha string) (*model.Commit, error) {
  var commit model.Commit

  if err := app.DB.Where(&model.Commit{ProjectID: projectID, Sha: sha}).FirstOrCreate(&commit).Error; err != nil {
    return nil, fmt.Errorf("error upserting Commit: %s", err.Error())
  }

  return &commit, nil
}