package model

import (
  "github.com/jinzhu/gorm"
  "github.com/sweettea-io/rest-api/internal/pkg/util/unique"
  "github.com/sweettea-io/rest-api/internal/pkg/util/str"
  "github.com/sweettea-io/rest-api/internal/pkg/util/enc"
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
  Name           string       `gorm:"default:null;not null"`
  Slug           string       `gorm:"default:null;not null;unique_index:deploy_slug"`
  DeploymentName string       `gorm:"default:null"`
  Public         bool         `gorm:"default:false"`
  Stage          string       `gorm:"default:null;not null"`
  Failed         bool         `gorm:"default:false"`
  LBHostname     string       `gorm:"type:varchar(240);default:null"`
  Hostname       string       `gorm:"type:varchar(240);default:null"`
  ClientID       string       `gorm:"type:varchar(240);default:null;not null"`
  ClientSecret   string       `gorm:"type:varchar(240);default:null;not null"`
  EnvVars        []EnvVar
}

// Assign Slug to Deploy before saving.
func (deploy *Deploy) BeforeSave(scope *gorm.Scope) error {
  scope.SetColumn("Slug", str.Slugify(deploy.Name))
  return nil
}

// Assign initial Stage, ClientID, & ClientSecret to Deploy before creation.
func (deploy *Deploy) BeforeCreate(scope *gorm.Scope) error {
  scope.SetColumn("Stage", BuildStages.Created)
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

func (deploy *Deploy) AsJSON() enc.JSON {
  mv := deploy.ModelVersion

  return enc.JSON{
    "name": deploy.Slug,
    "cluster": deploy.ApiCluster.Slug,
    "model": mv.Model.Slug,
    "modelVersion": mv.Version,
    "sha": deploy.Commit.Sha,
    "stage": deploy.Stage,
    "failed": deploy.Failed,
    "loadBalancer": deploy.LBHostname,
    "hostname": deploy.Hostname,
    "public": deploy.Public,
    "clientID": deploy.ClientID,
    "clientSecret": deploy.ClientSecret,
  }
}