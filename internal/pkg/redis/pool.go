package redis

import r "github.com/gomodule/redigo/redis"

// Create app's redis pool configured with app.Config's values
func NewPool(address string, password string, maxActive int, maxIdle int, wait bool) *r.Pool {
  pool := &r.Pool{
    // Max number of connections allocated by the pool at any given time (0 = unlimited)
    MaxActive: maxActive,

    // Max number of idle connections in the pool.
    MaxIdle: maxIdle,

    // Whether to wait for newly available connections if the pool is at its MaxActive limit.
    Wait: wait,

    Dial: func() (r.Conn, error) {
      // Connect to Redis address (hostname:port)
      c, err := r.Dial("tcp", address)

      if err != nil {
        panic(err)
      }

      // If Redis password is needed, use that.
      if password != "" {
        c.Do("AUTH", password)
      }

      return c, nil
    },
  }

  return pool
}