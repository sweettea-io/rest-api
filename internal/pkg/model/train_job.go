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

// TODO: implement this.
func (tj *TrainJob) ReadableStage() string {
  return "<train_job_stage>"
}

// TODO: implement this.
func (tj *TrainJob) ReadableCreatedAt() string {
  return "<train_job_created_at>"
}

func (tj *TrainJob) AsJSON() enc.JSON {
  return enc.JSON{
    "uid": tj.Uid,
    "commit": tj.Commit.Sha,
    "stage": tj.ReadableStage(),
    "failed": tj.Failed,
    "createdAt": tj.ReadableCreatedAt(),
    "model": enc.JSON{
      "name": tj.ModelVersion.Model.Slug,
      "version": tj.ModelVersion.Version,
    },
  }
}

func (tj *TrainJob) GetCommit() *Commit {
  return &tj.Commit
}

func (tj *TrainJob) GetUid() string {
  return tj.Uid
}