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
  gorm.Model
  Uid            string       `gorm:"type:varchar(240);default:null;not null;unique;index:train_job_uid"`
  Commit         Commit
  CommitID       uint         `gorm:"default:null;not null;unique_index:train_job_grouped_index"`
  ModelVersion   ModelVersion
  ModelVersionID uint         `gorm:"default:null;not null;unique_index:train_job_grouped_index"`
  Stage          uint         `gorm:"default:0"`
  Failed         bool         `gorm:"default:false"`
}

// Assign Uid to TrainJob before creation.
func (tj *TrainJob) BeforeCreate(scope *gorm.Scope) error {
  scope.SetColumn("Uid", unique.NewUid())
  return nil
}