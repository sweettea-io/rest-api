package model

import (
  "github.com/jinzhu/gorm"
  "github.com/sweettea-io/rest-api/internal/pkg/util/unique"
)

/*
  Commit <-- belongs_to

  has_one -- ModelVersion
*/
type TrainJob struct {
  // TODO: Model should be unique on commit_id+model_version_id (add a multi-col index too)
  gorm.Model
  Uid            string       `gorm:"type:varchar(240);default:null;not null;unique;index:train_job_uid"`
  Commit         Commit
  CommitID       uint         `gorm:"default:null;not null;index:train_job_commit_id"`
  ModelVersion   ModelVersion
  ModelVersionID uint         `gorm:"default:null;not null;index:train_job_model_version_id"`
  Stage          uint         `gorm:"default:0"`
  Failed         bool         `gorm:"default:false"`
}

// Assign Uid to TrainJob before creation.
func (tj *TrainJob) BeforeCreate(scope *gorm.Scope) error {
  scope.SetColumn("Uid", unique.NewUid())
  return nil
}