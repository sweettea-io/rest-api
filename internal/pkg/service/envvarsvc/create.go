package envvarsvc

import (
  "github.com/sweettea-io/rest-api/internal/pkg/model"
  "github.com/sweettea-io/rest-api/internal/app"
  "fmt"
)

func CreateFromMap(deployID uint, envs map[string]string) error {
  for key, val := range envs {
    if _, err := Create(deployID, key, val); err != nil {
      return err
    }
  }

  return nil
}

func Create(deployID uint, key string, val string) (*model.EnvVar, error) {
  //Create EnvVar model.
  envVar := model.EnvVar{
    DeployID: deployID,
    Key: key,
    Val: val,
  }

  // Create record.
  if err := app.DB.Create(&envVar).Error; err != nil {
    return nil, fmt.Errorf("error creating EnvVar: %s", err.Error())
  }

  return &envVar, nil
}

func UpsertFromMap(deployID uint, envs map[string]string) error {
  for key, val := range envs {
    if _, err := Upsert(deployID, key, val); err != nil {
      return err
    }
  }

  return nil
}

func Upsert(deployID uint, key string, val string) (*model.EnvVar, error) {
  var envVar model.EnvVar

  // Upsert EnvVar by DeployID+Key.
  if err := app.DB.
    Where(model.EnvVar{DeployID: deployID, Key: key}).
    Assign(model.EnvVar{Val: val}).
    FirstOrCreate(&envVar).Error; err != nil {
    return nil, fmt.Errorf("error upserting EnvVar: %s", err.Error())
  }

  return &envVar, nil
}