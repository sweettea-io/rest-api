package model

import (
  "time"
  "github.com/jinzhu/gorm"
  "github.com/sweettea-io/rest-api/internal/pkg/util/unique"
)

/*
  Model <-- belongs_to

  has_one -- TrainJob

  has_many --> Deploys
*/
type ModelVersion struct {
  // TODO: ModelVersion should be unique on model_id+version (add a multi-col index too)
  ID        uint       `gorm:"primary_key"`
  CreatedAt time.Time
  UpdatedAt time.Time
  DeletedAt *time.Time `sql:"index"`
  Model     Model
  ModelID   uint       `gorm:"default:null;not null;index:model_version_model_id"`
  Version   string     `gorm:"type:varchar(240);default:null;not null;index:model_version_version"`
  Deploys   []Deploy
}

// Assign Version to ModelVersion before creation.
func (mv *ModelVersion) BeforeCreate(scope *gorm.Scope) error {
  // TODO: Shorten version to 7 chars
  scope.SetColumn("Version", unique.NewUid())
  return nil
}