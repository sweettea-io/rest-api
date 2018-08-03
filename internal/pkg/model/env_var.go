package model

import (
  "github.com/jinzhu/gorm"
  "github.com/sweettea-io/rest-api/internal/pkg/util/unique"
)

/*
  Deploy <-- belongs_to
*/
type EnvVar struct {
  gorm.Model
  Uid        string `gorm:"type:varchar(240);default:null;not null;unique;index:env_var_uid"`
  Deploy     Deploy
  DeployID   uint   `gorm:"default:null;not null;unique_index:env_var_grouped_index"`
  Key        string `gorm:"type:varchar(360);default:null;not null;unique_index:env_var_grouped_index"`
  Val        string `gorm:"type:varchar(360);default:null"`
}

// Assign Uid to EnvVar before creation.
func (ev *EnvVar) BeforeCreate(scope *gorm.Scope) error {
  scope.SetColumn("Uid", unique.NewUid())
  return nil
}