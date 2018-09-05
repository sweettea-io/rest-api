package model

import (
  "github.com/jinzhu/gorm"
  "github.com/sweettea-io/rest-api/internal/pkg/util/enc"
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
  Stage          string       `gorm:"default:null;not null"`
  Failed         bool         `gorm:"default:false"`
}

// Assign initial stage before creation.
func (tj *TrainJob) BeforeCreate(scope *gorm.Scope) error {
  scope.SetColumn("Stage", BuildStages.Created)
  return nil
}

func (tj *TrainJob) AsJSON() enc.JSON {
  mv := tj.ModelVersion

  return enc.JSON{
    "model": mv.Model.Slug,
    "modelVersion": mv.Version,
    "sha": tj.Commit.Sha,
    "stage": tj.Stage,
    "failed": tj.Failed,
    "createdAt": tj.CreatedAt,
  }
}

func (tj *TrainJob) GetCommit() *Commit {
  return &tj.Commit
}

func (tj *TrainJob) GetUid() string {
  return tj.Uid
}