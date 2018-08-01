package model

import (
  "github.com/jinzhu/gorm"
  "github.com/sweettea-io/rest-api/internal/pkg/util/unique"
  "github.com/sweettea-io/rest-api/internal/pkg/model/buildable"
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
  Stage          string       `gorm:"default:null;not null"`
  Failed         bool         `gorm:"default:false"`
  LBHost         string       `gorm:"type:varchar(240);default:null"`
  ClientID       string       `gorm:"type:varchar(240);default:null;not null"`
  ClientSecret   string       `gorm:"type:varchar(240);default:null;not null"`
  EnvVars        []EnvVar
}

// Assign Uid, initial Stage, ClientID, & ClientSecret to Deploy before creation.
func (deploy *Deploy) BeforeCreate(scope *gorm.Scope) error {
  scope.SetColumn("Uid", unique.NewUid())
  scope.SetColumn("Stage", buildable.Created)
  scope.SetColumn("ClientID", unique.NewUid())
  scope.SetColumn("ClientSecret", unique.FreshSecret())
  return nil
}

func (deploy *Deploy) GetCommit() *Commit {
  return &deploy.Commit
}

func (deploy *Deploy) GetUid() string {
  return deploy.Uid
}