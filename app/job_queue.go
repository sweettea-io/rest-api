package app

import "github.com/sweettea-io/work"

var jobQueue *work.Enqueuer

// Create new job queue for passing jobs to our worker through redis
func CreateJobQueue() {
  jobQueue := work.NewEnqueuer(Config.JobQueueNsp, redisPool)
}