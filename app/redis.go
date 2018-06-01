package app

import "github.com/gomodule/redigo/redis"

var redisPool *redis.Pool

// Create app's redis pool configured with app.Config's values.
func CreateRedisPool() {
  redisPool = &redis.Pool{
    // Max number of connections allocated by the pool at a given time (0 = unlimited)
    MaxActive: Config.RedisPoolMaxActive,

    // Max number of idle connections in the pool.
    MaxIdle: Config.RedisPoolMaxIdle,

    // Wait for new connections to be available if pool is at its MaxActive limit?
    Wait: Config.RedisPoolWait,

    Dial: func() (redis.Conn, error) {
      return redis.Dial("tcp", Config.RedisUrl)
    },
  }
}