package app

import "github.com/gomodule/redigo/redis"

var RedisPool *redis.Pool

// Create app's redis pool configured with app.Config's values
func CreateRedisPool() {
  RedisPool = &redis.Pool{
    // Max number of connections allocated by the pool at any given time (0 = unlimited)
    MaxActive: Config.RedisPoolMaxActive,

    // Max number of idle connections in the pool.
    MaxIdle: Config.RedisPoolMaxIdle,

    // Whether to wait for newly available connections if the pool is at its MaxActive limit.
    Wait: Config.RedisPoolWait,

    Dial: func() (redis.Conn, error) {
      // Connect to Redis address (hostname:port)
      c, err := redis.Dial("tcp", Config.RedisAddress)

      if err != nil {
        panic(err)
      }

      // If Redis password is needed, use that.
      if Config.RedisPassword != "" {
        c.Do("AUTH", Config.RedisPassword)
      }

      return c, nil
    },
  }
}