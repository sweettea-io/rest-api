package deploysvc

import (
  "github.com/sweettea-io/rest-api/internal/pkg/model"
  "github.com/sweettea-io/rest-api/internal/app"
  "fmt"
)

func Upsert(commitID uint, modelVersionID uint, apiClusterID uint, uid string) (*model.Deploy, bool, error) {
  deployWithAttrs := model.Deploy{
    CommitID: commitID,
    ModelVersionID: modelVersionID,
    ApiClusterID: apiClusterID,
  }
  
  // Attempt to find Deploy.
  var deploy model.Deploy
  result := app.DB.Where(&deployWithAttrs).Find(&deploy)

  // If not found, create it.
  if result.RecordNotFound() {
    deploy = deployWithAttrs
    deploy.Uid = uid
    
    // Create record.
    if err := app.DB.Create(&deploy).Error; err != nil {
      return nil, true, fmt.Errorf("error creating Deploy: %s", err.Error())
    }

    return &deploy, true, nil
  }

  return &deploy, false, nil
}