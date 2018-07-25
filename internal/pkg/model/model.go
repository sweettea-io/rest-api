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
  // TODO: Model should be unique on project_id+slug (add a multi-col index too)
  gorm.Model
  Project       Project
  ProjectID     uint           `gorm:"default:null;not null;index:model_project_id"`
  Name          string         `gorm:"type:varchar(240);default:null;not null"`
  Slug          string         `gorm:"type:varchar(240);default:null;not null;index:model_slug"`
  FileExt       string         `gorm:"type:varchar(240);default:null"`
  ModelVersions []ModelVersion
}

// Assign Slug to Model before creation.
func (model *Model) BeforeCreate(scope *gorm.Scope) error {
  scope.SetColumn("Slug", str.Slugify(model.Name))
  return nil
}