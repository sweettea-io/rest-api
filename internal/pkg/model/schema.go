/*

SweetTea Database Schema

Models:

    User
    Session
    Project
    ProjectConfig
    Commit
    TrainJob
    Model
    ModelVersion
    Deploy
    Cluster
    EnvVar

Relationships:

    User|Session
        User --> has many --> sessions
        Session --> belongs to --> User

    Project|Commit
        Project --> has many --> commits
        Commit --> belongs to --> Project

    Project|ProjectConfig
        Project --> has one --> ProjectConfig
        ProjectConfig --> has one --> Project

    Commit|TrainJob
        Commit --> has many --> trainJobs
        TrainJob --> belongs to --> Commit

    Commit|Deploy
        Commit --> has many --> deploys
        Deploy --> belongs to --> Commit

    Model|ModelVersion
        Model --> has many --> modelVersions
        ModelVersion --> belongs to --> Model

    TrainJob|ModelVersion
        TrainJob --> has one --> ModelVersion
        ModelVersion --> has one --> TrainJob

    ModelVersion|Deploy
        ModelVersion --> has many --> deploys
        Deploy --> belongs to --> ModelVersion

    Cluster|Deploy
        Cluster --> has many --> deploys
        Deploy --> belongs to --> Cluster

    Deploy|EnvVar
        Deploy --> has many --> envVars
        EnvVar --> belongs to --> Deploy

*/
package model

import (
  "time"
  "github.com/jinzhu/gorm"
)

type User struct {
  gorm.Model
  Uid        string    `gorm:"type:varchar(240);default:null;not null;unique;index:user_uid"`
  Email      string    `gorm:"type:varchar(240);default:null;not null;unique;index:user_email"`
  HashedPw   string    `gorm:"type:varchar(240);default:null"`
  Admin      bool      `gorm:"default:false"`
  Sessions   []Session
  WithUid
}

type Session struct {
  gorm.Model
  User       User
  UserID     uint   `gorm:"default:null;not null;index:session_user_id"`
  Token      string `gorm:"type:varchar(240);default:null;not null;unique;index:session_token"`
}

type Project struct {
  gorm.Model
  Uid             string        `gorm:"type:varchar(240);default:null;not null;unique;index:project_uid"`
  Host            string        `gorm:"type:varchar(240);default:null;not null"`
  Nsp             string        `gorm:"type:varchar(360);default:null;not null;index:project_nsp"`
  ProjectConfig   ProjectConfig
  ProjectConfigID uint          `gorm:"default:null;not null;index:project_config_id"`
  Commits         []Commit
  WithUid
}

type ProjectConfig struct {
  gorm.Model
}

type Commit struct {
  gorm.Model
  Project    Project
  ProjectID  uint       `gorm:"default:null;not null;index:commit_project_id"`
  Sha        string     `gorm:"type:varchar(240);default:null;not null;unique;index:commit_sha"`
  Branch     string     `gorm:"type:varchar(240);default:null"`
  TrainJobs  []TrainJob
  Deploys    []Deploy
}

type TrainJob struct {
  gorm.Model
  Uid            string       `gorm:"type:varchar(240);default:null;not null;unique;index:train_job_uid"`
  Commit         Commit
  CommitID       uint         `gorm:"default:null;not null;index:train_job_commit_id"`
  ModelVersion   ModelVersion
  ModelVersionID uint         `gorm:"default:null;not null;index:train_job_model_version_id"`
  Stage          uint         `gorm:"default:0"`
  Failed         bool         `gorm:"default:false"`
  WithUid
}

type Model struct {
  gorm.Model
  Uid           string         `gorm:"type:varchar(240);default:null;not null;unique;index:model_uid"`
  Name          string         `gorm:"type:varchar(240);default:null;not null"`
  Slug          string         `gorm:"type:varchar(240);default:null;not null"`
  FileExt       string         `gorm:"type:varchar(240);default:null"`
  ModelVersions []ModelVersion
}

type ModelVersion struct {
  ID        uint       `gorm:"primary_key"`
  CreatedAt time.Time
  UpdatedAt time.Time
  DeletedAt *time.Time `sql:"index"`
  Model     Model
  ModelID   uint       `gorm:"default:null;not null;index:model_version_model_id"`
  Version   string     `gorm:"type:varchar(240);default:null;not null;unique;index:model_version_version"`
  Deploys   []Deploy
}

type Deploy struct {
  gorm.Model
  Uid            string       `gorm:"type:varchar(240);default:null;not null;unique;index:deploy_uid"`
  Commit         Commit
  CommitID       uint         `gorm:"default:null;not null;index:deploy_commit_id"`
  ModelVersion   ModelVersion
  ModelVersionID uint         `gorm:"default:null;not null;index:deploy_model_version_id"`
  Cluster        Cluster
  ClusterID      uint         `gorm:"default:null;not null;index:deploy_cluster_id"`
  Stage          uint         `gorm:"default:0"`
  Failed         bool         `gorm:"default:false"`
  LBHost         string       `gorm:"type:varchar(240);default:null"`
  ClientID       string       `gorm:"type:varchar(240);default:null;not null"`
  ClientSecret   string       `gorm:"type:varchar(240);default:null;not null"`
  EnvVars        []EnvVar
}

type Cluster struct {
  gorm.Model
  Uid        string   `gorm:"type:varchar(240);default:null;not null;unique;index:cluster_uid"`
  Name       string   `gorm:"type:varchar(360);default:null;not null"`
  Slug       string   `gorm:"type:varchar(360);default:null;not null;unique;index:cluster_slug"`
  Cloud      string   `gorm:"type:varchar(240);default:null;not null"`
  State      string   `gorm:"type:varchar(360);default:null"`
  Deploys    []Deploy
}

type EnvVar struct {
  gorm.Model
  Uid        string `gorm:"type:varchar(240);default:null;not null;unique;index:env_var_uid"`
  Deploy     Deploy
  DeployID   uint   `gorm:"default:null;not null;index:env_var_deploy_id"`
  Key        string `gorm:"type:varchar(360);default:null;not null"`
  Val        string `gorm:"type:varchar(360);default:null"`
  WithUid
}