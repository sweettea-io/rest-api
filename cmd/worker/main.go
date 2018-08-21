package main

import (
  "os"
  "os/signal"
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/app/worker/jobs"
  "github.com/sweettea-io/rest-api/internal/pkg/config"
  "github.com/sweettea-io/work"
)

func main() {
  // Initialize the app.
  app.Init(config.New())

  // Create new worker pool with empty job context.
  workerPool := work.NewWorkerPool(
    jobs.Context{},
    app.Config.WorkerCount,
    app.Config.JobQueueNsp,
    app.Redis,
  )

  // Add job middleware.
  workerPool.Middleware((*jobs.Context).LogJobStart)

  // Create job options that specify no retries for jobs that fail.
  noRetry := work.JobOptions{MaxFails: 1}

  // Assign handler functions to different jobs.
  workerPool.JobWithOptions(jobs.Names.CreateTrainJob, noRetry, (*jobs.Context).CreateTrainJob)
  workerPool.JobWithOptions(jobs.Names.CreateDeploy, noRetry, (*jobs.Context).CreateDeploy)
  workerPool.JobWithOptions(jobs.Names.BuildDeploy, noRetry, (*jobs.Context).BuildDeploy)
  workerPool.JobWithOptions(jobs.Names.TrainDeploy, noRetry, (*jobs.Context).TrainDeploy)
  workerPool.JobWithOptions(jobs.Names.ApiDeploy, noRetry, (*jobs.Context).ApiDeploy)
  workerPool.JobWithOptions(jobs.Names.UpdateDeploy, noRetry, (*jobs.Context).UpdateDeploy)
  workerPool.JobWithOptions(jobs.Names.ApiUpdate, noRetry, (*jobs.Context).ApiUpdate)
  workerPool.JobWithOptions(jobs.Names.PublicizeDeploy, noRetry, (*jobs.Context).PublicizeDeploy)

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

