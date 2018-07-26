package worker

import (
  "github.com/Sirupsen/logrus"
  "github.com/sweettea-io/work"
  "github.com/sweettea-io/rest-api/internal/app/worker/jobs"
  "github.com/jinzhu/gorm"
  "github.com/sweettea-io/rest-api/internal/pkg/config"
)

var JobContext *jobs.Context

func Init(cfg *config.Config, jobQueue *work.Enqueuer, db *gorm.DB, logger *logrus.Logger) {
  // Initialize job context.
  JobContext = &jobs.Context{
    Config: cfg,
    JobQueue: jobQueue,
    DB: db,
    Log: logger,
  }
}