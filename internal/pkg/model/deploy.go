package model

import (
  "github.com/jinzhu/gorm"
  "github.com/sweettea-io/rest-api/internal/pkg/util/unique"
)

/*
  Commit       <-- belongs_to
  ModelVersion <-- belongs_to
  ApiCluster   <-- belongs_to

  has_many --> EnvVars
*/
type Deploy struct {
  gorm.Model
  Uid            string       `gorm:"type:varchar(240);default:null;not null;unique;index:deploy_uid"`
  Commit         Commit
  CommitID       uint         `gorm:"default:null;not null;unique_index:deploy_grouped_index"`
  ModelVersion   ModelVersion
  ModelVersionID uint         `gorm:"default:null;not null;unique_index:deploy_grouped_index"`
  ApiCluster     ApiCluster
  ApiClusterID   uint         `gorm:"default:null;not null;unique_index:deploy_grouped_index"`
  Stage          uint         `gorm:"default:0"`
  Failed         bool         `gorm:"default:false"`
  LBHost         string       `gorm:"type:varchar(240);default:null"`
  ClientID       string       `gorm:"type:varchar(240);default:null;not null"`
  ClientSecret   string       `gorm:"type:varchar(240);default:null;not null"`
  EnvVars        []EnvVar
}

// Assign Uid, ClientID, & ClientSecret to Deploy before creation.
func (deploy *Deploy) BeforeCreate(scope *gorm.Scope) error {
  scope.SetColumn("Uid", unique.NewUid())
  scope.SetColumn("ClientID", unique.NewUid())
  scope.SetColumn("ClientSecret", unique.FreshSecret())
  return nil
}
