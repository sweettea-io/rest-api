package jobs

import (
  "github.com/Sirupsen/logrus"
  "github.com/sweettea-io/work"
  "github.com/sweettea-io/rest-api/internal/pkg/config"
  "github.com/jinzhu/gorm"
)

// Context represents job context for each job enqueued.
// Job handlers are attached as functions to this struct.
type Context struct {
  Config *config.Config
  JobQueue *work.Enqueuer
  DB *gorm.DB
  Log *logrus.Logger
}