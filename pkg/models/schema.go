// Models:
//
//    User
//    Session
//    Company
//    Cluster
//    Bucket
//    Project
//    Dataset
//    Env
//    Commit
//    Deploy
//
// Relationships:
//
//    User|Session
//        User --> has_many --> sessions
//        Session --> belongs_to --> User
//
//    Company|Cluster
//        Company --> has_one --> Cluster
//        Cluster --> has_one --> Company
//
//    Company|Bucket
//        Company --> has_one --> Bucket
//        Bucket --> has_one --> Company
//
//    Company|Project
//        Company --> has_many --> projects
//        Project --> belongs_to --> Company
//
//    Project|Dataset
//        Project --> has_many --> datasets
//        Dataset --> belongs_to --> Project
//
//    Project|Env
//        Project --> has_many --> envs
//        Env --> belongs_to --> Project
//
//    Project|Commit
//        Project --> has_many --> commits
//        Commit --> belongs_to --> Project
//
//    Commit|Deploys
//        Commit --> has_many --> deploys
//        Deploys --> belongs_to --> Commit
//
package models

import (
  "time"
  "github.com/jinzhu/gorm"
  "github.com/sweettea-io/rest-api/pkg/utils"
)

type WithUid struct {}

type User struct {
  gorm.Model
  Uid         string    `gorm:"type:varchar(240);default:null;unique;not null;index:user_uid"`
  Email       string    `gorm:"type:varchar(240);default:null;unique;not null;index:user_email"`
  HashedPw    string    `gorm:"type:varchar(240);default:null"`
  Admin       bool      `gorm:"default:false"`
  Sessions    []Session
  IsDestroyed bool      `gorm:"default:false"`
  WithUid
}

type Session struct {
  gorm.Model
  User       User
  UserID     int    `gorm:"default:null;not null;index:session_user_id"`
  Token      string `gorm:"type:varchar(240);default:null;unique;not null;index:session_token"`
}

type Company struct {
  gorm.Model
  Uid         string  `gorm:"type:varchar(240);default:null;unique;not null;index:company_uid"`
  Name        string  `gorm:"type:varchar(240);default:null;not null"`
  Slug        string  `gorm:"type:varchar(240);default:null;unique;not null;index:company_slug"`
  Cluster     Cluster
  ClusterID   int     `gorm:"default:null;index:company_cluster_id"`
  Bucket      Bucket
  BucketID    int     `gorm:"default:null;index:company_bucket_id"`
  Projects       []Project
  IsDestroyed bool    `gorm:"default:false"`
  WithUid
}

type Cluster struct {
  gorm.Model
  Uid          string      `gorm:"type:varchar(240);default:null;unique;not null;index:cluster_uid"`
  Name         string      `gorm:"type:varchar(360);default:null;not null"`
  NsAddresses  utils.JsonB `gorm:"type:jsonb;not null;default:'{}'::jsonb"`
  Zones        utils.JsonB `gorm:"type:jsonb;not null;default:'{}'::jsonb"`
  HostedZoneID string      `gorm:"type:varchar(240);default:null"`
  MasterType   string      `gorm:"type:varchar(240);default:null"`
  NodeType     string      `gorm:"type:varchar(240);default:null"`
  Image        string      `gorm:"type:varchar(240);default:null"`
  Validated    bool        `gorm:"default:false"`
  IsDestroyed  bool        `gorm:"default:false"`
  WithUid
}

type Bucket struct {
  gorm.Model
  Name        string `gorm:"type:varchar(240);default:null"`
  IsDestroyed bool   `gorm:"default:false"`
}

type Project struct {
  gorm.Model
  Uid              string    `gorm:"type:varchar(240);default:null;unique;not null;index:project_uid"`
  Name             string    `gorm:"type:varchar(240);default:null;not null"`
  Slug             string    `gorm:"type:varchar(240);default:null;unique;not null;index:project_slug"`
  Company          Company
  CompanyID        int       `gorm:"default:null;not null;index:project_company_id"`
  Datasets         []Dataset
  Envs             []Env
  Commits          []Commit
  Elb              string    `gorm:"type:varchar(240);default:null"`
  Domain           string    `gorm:"type:varchar(360);default:null"`
  DeployName       string    `gorm:"type:varchar(360);default:null"`
  ClientID         string    `gorm:"type:varchar(240);default:null"`
  ClientSecret     string    `gorm:"type:varchar(240);default:null"`
  ModelExt         string    `gorm:"type:varchar(240);default:null"`
  InternalMsgToken string    `gorm:"type:varchar(240);default:null"`
  IsDestroyed      bool      `gorm:"default:false"`
  WithUid
}

type Dataset struct {
  gorm.Model
  Uid                  string `gorm:"type:varchar(240);default:null;unique;not null;index:dataset_uid"`
  Name                 string `gorm:"type:varchar(240);default:null;not null"`
  Slug                 string `gorm:"type:varchar(240);default:null;unique;not null;index:dataset_slug"`
  Project                 Project
  ProjectID               int    `gorm:"default:null;not null;index:dataset_project_id"`
  RetrainStepSize      int    `gorm:"default:0"`
  LastTrainRecordCount int    `gorm:"default:0"`
  IsDestroyed          bool   `gorm:"default:false"`
  WithUid
}

type Env struct {
  gorm.Model
  Uid        string `gorm:"type:varchar(240);default:null;unique;not null;index:env_uid"`
  Project       Project
  ProjectID     int    `gorm:"default:null;not null;index:env_project_id"`
  Name       string `gorm:"type:varchar(240);default:null;not null"`
  Value      string `gorm:"type:varchar(240);default:null;not null"`
  ForCluster string `gorm:"type:varchar(240);default:null;not null"`
  WithUid
}

type Commit struct {
  gorm.Model
  Project        Project
  ProjectID      int          `gorm:"default:null;not null;index:commit_project_id"`
  Deployss []Deploys
  Sha         string       `gorm:"type:varchar(240);default:null;unique;not null;index:commit_sha"`
  Message     string       `gorm:"type:varchar(240);default:null"`
  Author      string       `gorm:"type:varchar(240);default:null"`
  Branch      string       `gorm:"type:varchar(240);default:null"`
}

type Deploys struct {
  gorm.Model
  Uid              string    `gorm:"type:varchar(240);default:null;unique;not null;index:deploy_uid"`
  Commit           Commit
  CommitID         int       `gorm:"default:null;not null;index:deploy_commit_id"`
  Status           string    `gorm:"type:varchar(240);default:null"`
  TrainTriggeredBy string    `gorm:"type:varchar(240);default:null"`
  ServeTriggeredBy string    `gorm:"type:varchar(240);default:null"`
  Intent           string    `gorm:"type:varchar(240);default:null"`
  IntentUpdatedAt  time.Time
  Failed           bool      `gorm:"default:false"`
  WithUid
}

// -------- Model Creation-related Hooks ----------

// Assign Uid to model before creation.
func (w *WithUid) BeforeCreate(scope *gorm.Scope) error {
  scope.SetColumn("Uid", utils.NewUid())
  return nil
}

// Assign newly minted secret to Session token before creation.
func (session *Session) BeforeCreate(scope *gorm.Scope) error {
  scope.SetColumn("Token", utils.FreshSecret())
  return nil
}

// Assign Uid & Slug to Company before creation.
func (company *Company) BeforeCreate(scope *gorm.Scope) error {
  scope.SetColumn("Uid", utils.NewUid())
  scope.SetColumn("Slug", utils.Slugify(company.Name))
  return nil
}

// Assign Uid & Slug to Project before creation.
func (project *Project) BeforeCreate(scope *gorm.Scope) error {
  scope.SetColumn("Uid", utils.NewUid())
  scope.SetColumn("Slug", utils.Slugify(project.Name))
  return nil
}

// Assign Uid & Slug to Dataset before creation.
func (dataset *Dataset) BeforeCreate(scope *gorm.Scope) error {
  scope.SetColumn("Uid", utils.NewUid())
  scope.SetColumn("Slug", utils.Slugify(dataset.Name))
  return nil
}