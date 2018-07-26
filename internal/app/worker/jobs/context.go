package jobs

import (
  "github.com/Sirupsen/logrus"
  "github.com/sweettea-io/work"
)

// Context represents job context for each job enqueued.
// Job handlers are attached as functions to this struct.
type Context struct {
  Log *logrus.Logger
  JobQueue *work.Enqueuer
}