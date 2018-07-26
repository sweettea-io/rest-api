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
  ID        uint       `gorm:"primary_key"`
  CreatedAt time.Time
  UpdatedAt time.Time
  DeletedAt *time.Time `sql:"index"`
  Model     Model
  ModelID   uint       `gorm:"default:null;not null;unique_index:model_version_grouped_index"`
  Version   string     `gorm:"type:varchar(240);default:null;not null;unique_index:model_version_grouped_index"`
  Deploys   []Deploy
}

// Assign Version to ModelVersion before creation.
func (mv *ModelVersion) BeforeCreate(scope *gorm.Scope) error {
  scope.SetColumn("Version", unique.NewUid()[:8])
  return nil
}