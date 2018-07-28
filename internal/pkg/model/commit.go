package model

import "github.com/jinzhu/gorm"

/*
  Project <-- belongs_to

  has_many --> TrainJobs
  has_many --> Deploys
*/
type Commit struct {
  gorm.Model
  Project    Project
  ProjectID  uint       `gorm:"default:null;not null;index:commit_project_id"`
  Sha        string     `gorm:"type:varchar(240);default:null;not null;unique;index:commit_sha"`
  TrainJobs  []TrainJob
  Deploys    []Deploy
}