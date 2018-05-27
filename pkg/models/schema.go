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
  "github.com/jinzhu/gorm"
  "time"
)

type Destroyable struct {
  IsDestroyed bool `gorm:"default:false"`
}

type WithUid struct {
  Uid string `gorm:"type:varchar(240);unique;not null;unique_index"`
}

type WithSlug struct {
  Slug string `gorm:"type:varchar(240);unique;not null;unique_index"`
}

type User struct {
  gorm.Model
  WithUid
  Email       string    `gorm:"type:varchar(240);unique;not null;unique_index"`
  HashedPw    string    `gorm:"type:varchar(240)"`
  Sessions    []Session
  Destroyable
}

type Session struct {
  gorm.Model
  User       User
  UserID     int    `gorm:"not null;unique_index"`
  Token      string `gorm:"type:varchar(240)"`
}

type Team struct {
  gorm.Model
  WithUid
  Name        string  `gorm:"type:varchar(240);not null"`
  WithSlug
  Cluster     Cluster
  ClusterID   int     `gorm:"unique_index"`
  Bucket      Bucket
  BucketID    int     `gorm:"unique_index"`
  Repos       []Repo
  Destroyable
}

type Cluster struct {
  gorm.Model
  WithUid
  Name         string `gorm:"type:varchar(360);not null"`
  NsAddresses  JSON   `gorm:"type:jsonb;not null;default:'{}'::jsonb"`
  Zones        JSON   `gorm:"type:jsonb;not null;default:'{}'::jsonb"`
  HostedZoneID string `gorm:"type:varchar(240)"`
  MasterType   string `gorm:"type:varchar(240)"`
  NodeType     string `gorm:"type:varchar(240)"`
  Image        string `gorm:"type:varchar(240)"`
  Validated    bool   `gorm:"default:false"`
  Destroyable
}

type Bucket struct {
  gorm.Model
  Name        string `gorm:"type:varchar(240)"`
  Destroyable
}

type Repo struct {
  gorm.Model
  WithUid
  Name             string    `gorm:"type:varchar(240);not null"`
  WithSlug
  Team             Team
  TeamID           int       `gorm:"not null;unique_index"`
  Datasets         []Dataset
  Envs             []Env
  Commits          []Commit
  Elb              string    `gorm:"type:varchar(240)"`
  Domain           string    `gorm:"type:varchar(360)"`
  ImageRepoOwner   string    `gorm:"type:varchar(240)"`
  DeployName       string    `gorm:"type:varchar(360)"`
  ClientID         string    `gorm:"type:varchar(240)"`
  ClientSecret     string    `gorm:"type:varchar(240)"`
  ModelExt         string    `gorm:"type:varchar(240)"`
  InternalMsgToken string    `gorm:"type:varchar(240)"`
  Destroyable
}

type Dataset struct {
  gorm.Model
  WithUid
  Name                 string `gorm:"type:varchar(240);not null"`
  WithSlug
  Repo                 Repo
  RepoID               int    `gorm:"not null;unique_index"`
  RetrainStepSize      int
  LastTrainRecordCount int
  Destroyable
}

type Env struct {
  gorm.Model
  WithUid
  Repo       Repo
  RepoID     int    `gorm:"not null;unique_index"`
  Name       string `gorm:"type:varchar(240);not null"`
  Value      string `gorm:"type:varchar(240);not null"`
  ForCluster string `gorm:"type:varchar(240);not null"`
}

type Commit struct {
  gorm.Model
  Repo        Repo
  RepoID      int          `gorm:"not null;unique_index"`
  Deployments []Deployment
  Sha         string       `gorm:"type:varchar(240);unique;not null;unique_index"`
  Message     string       `gorm:"type:varchar(240)"`
  Author      string       `gorm:"type:varchar(240)"`
  Branch      string       `gorm:"type:varchar(240)"`
}

type Deployment struct {
  gorm.Model
  WithUid
  Commit           Commit
  CommitID         int       `gorm:"not null;unique_index"`
  Status           string    `gorm:"type:varchar(240)"`
  TrainTriggeredBy string    `gorm:"type:varchar(240)"`
  ServeTriggeredBy string    `gorm:"type:varchar(240)"`
  Intent           string    `gorm:"type:varchar(240)"`
  IntentUpdatedAt  time.Time
  Failed           bool      `gorm:"default:false"`
}