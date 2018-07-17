package app

import (
  "github.com/jinzhu/gorm"
  "github.com/Sirupsen/logrus"
  "github.com/sweettea-io/rest-api/internal/pkg/config"
  "github.com/sweettea-io/rest-api/internal/pkg/db"
  "github.com/sweettea-io/rest-api/internal/pkg/logger"
  "github.com/sweettea-io/rest-api/internal/pkg/redis"
  "github.com/sweettea-io/work"
  r "github.com/gomodule/redigo/redis"
)

var Config *config.Config
var Redis *r.Pool
var JobQueue *work.Enqueuer
var DB *gorm.DB
var Log *logrus.Logger

func Init(cfg *config.Config) {
  // Set global config.
  Config = cfg

  // Create redis pool.
  Redis = redis.NewPool(
    Config.RedisAddress,
    Config.RedisPassword,
    Config.RedisPoolMaxActive,
    Config.RedisPoolMaxIdle,
    Config.RedisPoolWait,
  )

  // Create worker job queue.
  JobQueue = work.NewEnqueuer(Config.JobQueueNsp, Redis)

  // Create postgres database connection.
  DB = db.NewConnection(Config.DatabaseUrl, Config.Debug)

  // Create app logger.
  Log = logger.NewLogger()
}