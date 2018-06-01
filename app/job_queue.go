package app

import "github.com/sweettea-io/work"

var JobQueue *work.Enqueuer

// Create new job queue for passing jobs to our worker through redis
func CreateJobQueue() {
  JobQueue = work.NewEnqueuer(Config.JobQueueNsp, RedisPool)
}