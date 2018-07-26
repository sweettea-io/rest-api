package main

import (
  "os"
  "os/signal"
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/app/worker"
  "github.com/sweettea-io/rest-api/internal/app/worker/jobs"
  "github.com/sweettea-io/rest-api/internal/pkg/config"
  "github.com/sweettea-io/work"
)

func main() {
  // Initialize the app.
  app.Init(config.New())

  // Initialize the worker.
  worker.Init(
    app.Config,
    app.JobQueue,
    app.DB,
    app.Log,
  )

  // Create new worker pool with job context.
  workerPool := work.NewWorkerPool(
    worker.JobContext,
    app.Config.WorkerCount,
    app.Config.JobQueueNsp,
    app.Redis,
  )

  // Add job middleware.
  workerPool.Middleware((*jobs.Context).LogJobStart)

  // Assign handler functions to different jobs.
  workerPool.Job(jobs.Names.CreateTrainJob, (*jobs.Context).CreateTrainJob)

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

