package jobs

import (
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/work"
)

// Simple job middleware used to log whenever a new job is started.
func (c *Context) LogJobStart(job *work.Job, next work.NextMiddlewareFunc) error {
  app.Log.Infof("Starting job %s...\n", job.Name)
  return next()
}