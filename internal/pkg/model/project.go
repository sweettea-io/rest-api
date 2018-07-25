package model

import (
  "github.com/jinzhu/gorm"
  "github.com/sweettea-io/rest-api/internal/pkg/util/unique"
)

/*
  has_one -- ProjectConfig

  has_many --> Commits
  has_many --> Models
*/
type Project struct {
  gorm.Model
  Uid             string        `gorm:"type:varchar(240);default:null;not null;unique;index:project_uid"`
  Host            string        `gorm:"type:varchar(240);default:null;not null"`
  Nsp             string        `gorm:"type:varchar(360);default:null;not null;unique;index:project_nsp"`
  ProjectConfig   ProjectConfig
  ProjectConfigID uint          `gorm:"default:null;not null;index:project_config_id"`
  Commits         []Commit
  Models          []Model
}

// Assign Uid to Project before creation.
func (project *Project) BeforeCreate(scope *gorm.Scope) error {
  scope.SetColumn("Uid", unique.NewUid())
  return nil
}