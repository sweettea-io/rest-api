package main

import (
  "os"
  "os/signal"
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/Sirupsen/logrus"
  "github.com/sweettea-io/rest-api/internal/pkg/config"
  "github.com/sweettea-io/work"
)

// Job context with methods (defined below) representing middleware used for each job.
type Context struct {
  logger *logrus.Logger
}

// Simple job middleware used to log whenever a new job is started.
func (c *Context) LogJobStart(job *work.Job, next work.NextMiddlewareFunc) error {
  c.logger.Infof("Starting job: %s", job.Name)
  return next()
}

func main() {
  // Initialize the app.
  app.Init(config.New())

  // Create job context.
  context := Context{logger: app.Log}

  // Create new worker pool with job context.
  workerPool := work.NewWorkerPool(
    context,
    app.Config.WorkerCount,
    app.Config.JobQueueNsp,
    app.Redis,
  )

  // Add middleware that will be executed for each job.
  workerPool.Middleware((*Context).LogJobStart)

  // Assign handler functions to different jobs (by name).
  // TODO

  // Start processing jobs.
  workerPool.Start()

  app.Log.Infof("Starting worker pool with %v workers...", app.Config.WorkerCount)

  // Wait for signal to quit.
  signalChan := make(chan os.Signal, 1)
  signal.Notify(signalChan, os.Interrupt, os.Kill)
  <-signalChan

  // Stop worker pool.
  workerPool.Stop()
}

