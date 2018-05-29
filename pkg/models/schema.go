// Models:
//
//    User
//    Session
//    Team
//    Cluster
//    Bucket
//    Repo
//    Dataset
//    Env
//    Commit
//    Deployment
//
// Relationships:
//
//    User|Session
//        User --> has_many --> sessions
//        Session --> belongs_to --> User
//
//    Team|Cluster
//        Team --> has_one --> Cluster
//        Cluster --> has_one --> Team
//
//    Team|Bucket
//        Team --> has_one --> Bucket
//        Bucket --> has_one --> Team
//
//    Team|Repo
//        Team --> has_many --> repos
//        Repo --> belongs_to --> Team
//
//    Repo|Dataset
//        Repo --> has_many --> datasets
//        Dataset --> belongs_to --> Repo
//
//    Repo|Env
//        Repo --> has_many --> envs
//        Env --> belongs_to --> Repo
//
//    Repo|Commit
//        Repo --> has_many --> commits
//        Commit --> belongs_to --> Repo
//
//    Commit|Deployment
//        Commit --> has_many --> deployments
//        Deployment --> belongs_to --> Commit
//
package models

import (
  "time"
  "github.com/jinzhu/gorm"
  "github.com/sweettea/rest-api/pkg/utils"
)

type WithUid struct {}

type User struct {
  gorm.Model
  Uid         string    `gorm:"type:varchar(240);default:null;unique;not null;index:user_uid"`
  Email       string    `gorm:"type:varchar(240);default:null;unique;not null;index:user_email"`
  HashedPw    string    `gorm:"type:varchar(240);default:null"`
  Sessions    []Session
  IsDestroyed bool      `gorm:"default:false"`
  WithUid
}

type Session struct {
  gorm.Model
  User       User
  UserID     int    `gorm:"default:null;not null;index:session_user_id"`
  Token      string `gorm:"type:varchar(240);default:null"`
}

type Team struct {
  gorm.Model
  Uid         string  `gorm:"type:varchar(240);default:null;unique;not null;index:team_uid"`
  Name        string  `gorm:"type:varchar(240);default:null;not null"`
  Slug        string  `gorm:"type:varchar(240);default:null;unique;not null;index:team_slug"`
  Cluster     Cluster
  ClusterID   int     `gorm:"default:null;index:team_cluster_id"`
  Bucket      Bucket
  BucketID    int     `gorm:"default:null;index:team_bucket_id"`
  Repos       []Repo
  IsDestroyed bool    `gorm:"default:false"`
  WithUid
}

type Cluster struct {
  gorm.Model
  Uid          string     `gorm:"type:varchar(240);default:null;unique;not null;index:cluster_uid"`
  Name         string     `gorm:"type:varchar(360);default:null;not null"`
  NsAddresses  utils.JSON `gorm:"type:jsonb;not null;default:'{}'::jsonb"`
  Zones        utils.JSON `gorm:"type:jsonb;not null;default:'{}'::jsonb"`
  HostedZoneID string     `gorm:"type:varchar(240);default:null"`
  MasterType   string     `gorm:"type:varchar(240);default:null"`
  NodeType     string     `gorm:"type:varchar(240);default:null"`
  Image        string     `gorm:"type:varchar(240);default:null"`
  Validated    bool       `gorm:"default:false"`
  IsDestroyed  bool       `gorm:"default:false"`
  WithUid
}

type Bucket struct {
  gorm.Model
  Name        string `gorm:"type:varchar(240);default:null"`
  IsDestroyed bool   `gorm:"default:false"`
}

type Repo struct {
  gorm.Model
  Uid              string    `gorm:"type:varchar(240);default:null;unique;not null;index:repo_uid"`
  Name             string    `gorm:"type:varchar(240);default:null;not null"`
  Slug             string    `gorm:"type:varchar(240);default:null;unique;not null;index:repo_slug"`
  Team             Team
  TeamID           int       `gorm:"default:null;not null;index:repo_team_id"`
  Datasets         []Dataset
  Envs             []Env
  Commits          []Commit
  Elb              string    `gorm:"type:varchar(240);default:null"`
  Domain           string    `gorm:"type:varchar(360);default:null"`
  ImageRepoOwner   string    `gorm:"type:varchar(240);default:null"`
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
  Repo                 Repo
  RepoID               int    `gorm:"default:null;not null;index:dataset_repo_id"`
  RetrainStepSize      int    `gorm:"default:0"`
  LastTrainRecordCount int    `gorm:"default:0"`
  IsDestroyed          bool   `gorm:"default:false"`
  WithUid
}

type Env struct {
  gorm.Model
  Uid        string `gorm:"type:varchar(240);default:null;unique;not null;index:env_uid"`
  Repo       Repo
  RepoID     int    `gorm:"default:null;not null;index:env_repo_id"`
  Name       string `gorm:"type:varchar(240);default:null;not null"`
  Value      string `gorm:"type:varchar(240);default:null;not null"`
  ForCluster string `gorm:"type:varchar(240);default:null;not null"`
  WithUid
}

type Commit struct {
  gorm.Model
  Repo        Repo
  RepoID      int          `gorm:"default:null;not null;index:commit_repo_id"`
  Deployments []Deployment
  Sha         string       `gorm:"type:varchar(240);default:null;unique;not null;index:commit_sha"`
  Message     string       `gorm:"type:varchar(240);default:null"`
  Author      string       `gorm:"type:varchar(240);default:null"`
  Branch      string       `gorm:"type:varchar(240);default:null"`
}

type Deployment struct {
  gorm.Model
  Uid              string    `gorm:"type:varchar(240);default:null;unique;not null;index:deployment_uid"`
  Commit           Commit
  CommitID         int       `gorm:"default:null;not null;index:deployment_commit_id"`
  Status           string    `gorm:"type:varchar(240);default:null"`
  TrainTriggeredBy string    `gorm:"type:varchar(240);default:null"`
  ServeTriggeredBy string    `gorm:"type:varchar(240);default:null"`
  Intent           string    `gorm:"type:varchar(240);default:null"`
  IntentUpdatedAt  time.Time
  Failed           bool      `gorm:"default:false"`
  WithUid
}

// -------- Model Creation-related Hooks ----------

// Hook generating a uid for a model that needs one.
func (w *WithUid) BeforeCreate(scope *gorm.Scope) error {
  scope.SetColumn("Uid", utils.NewUid())
  return nil
}

// Assign Uid & Slug to Team before creation.
func (team *Team) BeforeCreate(scope *gorm.Scope) error {
  scope.SetColumn("Uid", utils.NewUid())
  scope.SetColumn("Slug", utils.Slugify(team.Name))
  return nil
}

// Assign Uid & Slug to Repo before creation.
func (repo *Repo) BeforeCreate(scope *gorm.Scope) error {
  scope.SetColumn("Uid", utils.NewUid())
  scope.SetColumn("Slug", utils.Slugify(repo.Name))
  return nil
}

// Assign Uid & Slug to Dataset before creation.
func (dataset *Dataset) BeforeCreate(scope *gorm.Scope) error {
  scope.SetColumn("Uid", utils.NewUid())
  scope.SetColumn("Slug", utils.Slugify(dataset.Name))
  return nil
}