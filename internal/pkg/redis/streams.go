package redis

import (
  r "github.com/gomodule/redigo/redis"
  "net"
  "time"
)

const DefaultReadTimeout = 30 // seconds
const DefaultXStartTs = "0.0" // very beginning of stream

func XRead(conn *r.Conn, stream string, startTs string, timeout int) (interface{}, error) {
  reply, err := r.DoWithTimeout(
    *conn,
    time.Duration(timeout) * time.Second,
    "XREAD",
    "BLOCK",
    timeout * 1000,
    "STREAMS",
    stream,
    startTs,
  )

  // If no error, return successfully with reply.
  if err == nil {
    return reply, nil
  }

  // If timeout error, return nil for both.
  if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
    return nil, nil
  }

  // Lastly, an unexpected error occurred (not a timeout error), so return that.
  return nil, err
}
