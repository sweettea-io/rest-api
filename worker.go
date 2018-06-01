package main

import (
  "github.com/sweettea/rest-api/app"
  "github.com/Sirupsen/logrus"
  "github.com/sweettea-io/work"
  "os"
  "os/signal"
)

// Job context with methods (defined below) representing middleware used for each job.
type Context struct {
  logger *logrus.Logger
}

// Simple job middleware used to log whenever a new job is started.
func (c *Context) Log(job *work.Job, next work.NextMiddlewareFunc) error {
  c.logger.Infof("Starting job: %s", job.Name)
  return next()
}

func main() {
  // Load app config.
  app.LoadConfig()

  // Create redis pool.
  app.CreateRedisPool()

  // Create logger.
  logger := logrus.New()

  // Create job context.
  context := Context{
    logger: logger,
  }

  // Create new worker pool with job context.
  workerPool := work.NewWorkerPool(
    context,
    app.Config.WorkerCount,
    app.Config.JobQueueNsp,
    app.RedisPool,
  )

  // Add middleware that will be executed for each job.
  workerPool.Middleware((*Context).Log)

  // Assign handler functions to different jobs (by name).
  // TODO

  // Start processing jobs.
  workerPool.Start()

  logger.Infof("Starting worker pool with %v workers...", app.Config.WorkerCount)

  // Wait for signal to quit.
  signalChan := make(chan os.Signal, 1)
  signal.Notify(signalChan, os.Interrupt, os.Kill)
  <-signalChan

  // Stop worker pool.
  workerPool.Stop()
}

