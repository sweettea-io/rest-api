package model

import (
  "github.com/jinzhu/gorm"
  "github.com/sweettea-io/rest-api/internal/pkg/util/unique"
  "strings"
  "github.com/sweettea-io/rest-api/internal/pkg/util/enc"
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

// Uppercase Key before save.
func (ev *EnvVar) BeforeSave(scope *gorm.Scope) error {
  scope.SetColumn("Key", strings.ToUpper(ev.Key))
  return nil
}

// Assign Uid to EnvVar before creation.
func (ev *EnvVar) BeforeCreate(scope *gorm.Scope) error {
  scope.SetColumn("Uid", unique.NewUid())
  return nil
}

func (ev *EnvVar) AsJSON() enc.JSON {
  return enc.JSON{
    "key": ev.Key,
    "value": ev.Val,
  }
}