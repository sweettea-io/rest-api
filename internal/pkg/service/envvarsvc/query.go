package envvarsvc

import (
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/pkg/model"
)

func GetMap(deployID uint) map[string]string {
  var envVars []model.EnvVar
  app.DB.Where(&model.EnvVar{DeployID: deployID}).Find(&envVars)

  envsMap := map[string]string{}
  for _, envVar := range envVars {
    envsMap[envVar.Key] = envVar.Val
  }

  return envsMap
}