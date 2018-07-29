package timeutil

import "time"

func MSSinceEpoch() int64 {
  return time.Now().UnixNano() / 1000000
}