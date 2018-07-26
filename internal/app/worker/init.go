package worker

import (
  "github.com/Sirupsen/logrus"
  "github.com/sweettea-io/work"
  "github.com/sweettea-io/rest-api/internal/app/worker/jobs"
)

var JobContext *jobs.Context

func Init(logger *logrus.Logger, queue *work.Enqueuer) {
  // Initialize job context.
  JobContext = &jobs.Context{
    Log: logger,
    JobQueue: queue,
  }
}