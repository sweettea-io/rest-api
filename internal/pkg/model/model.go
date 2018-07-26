package model

import (
  "github.com/jinzhu/gorm"
  "github.com/sweettea-io/rest-api/internal/pkg/util/str"
)

/*
  Project <-- belongs_to

  has_many --> ModelVersions
*/
type Model struct {
  gorm.Model
  Project       Project
  ProjectID     uint           `gorm:"default:null;not null;unique_index:model_grouped_index"`
  Name          string         `gorm:"type:varchar(240);default:null;not null"`
  Slug          string         `gorm:"type:varchar(240);default:null;not null;unique_index:model_grouped_index"`
  FileExt       string         `gorm:"type:varchar(240);default:null"`
  ModelVersions []ModelVersion
}

// Assign Slug to Model before creation.
func (model *Model) BeforeSave(scope *gorm.Scope) error {
  scope.SetColumn("Slug", str.Slugify(model.Name))
  return nil
}