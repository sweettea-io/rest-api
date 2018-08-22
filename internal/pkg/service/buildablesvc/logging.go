package buildablesvc

import (
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/pkg/redis"
)

type BuildableLog struct {
  Msg       string
  Timestamp string
  Completed bool
  Failed    bool
  Error     error
}

var StreamKeepAlive = BuildableLog{Msg: "..."}

func UnmarshalLog(raw interface{}) BuildableLog {
  log := BuildableLog{}

  // Set Completed if completed
  // Set Failed if level is 'error'

  return log
}

func NewErrorLog(err error) BuildableLog {
  return BuildableLog{
    Msg: "unexpected log stream error",
    Error: err,
  }
}

type LogStreamer struct {
  StreamKey string
  Channel   <-chan BuildableLog
}

func NewLogStreamer(buildableUid string) *LogStreamer {
  return &LogStreamer{
    StreamKey: buildableUid,
    Channel: make(chan BuildableLog),
  }
}

func (ls *LogStreamer) GetChannel() <-chan BuildableLog {
  return ls.Channel
}

func (ls *LogStreamer) Watch() {
  done := false
  readFromTs := redis.DefaultXStartTs // start reading from beginning of stream.

  // Get new connection from Redis pool.
  conn := app.Redis.Get()
  defer conn.Close()

  // Read from Redis stream until complete.
  for !done {
    reply, err := redis.XRead(&conn, ls.StreamKey, readFromTs, redis.DefaultReadTimeout)

    // Pipe any unexpected errors and return (if those occur).
    if err != nil {
      ls.Channel <- NewErrorLog(err)
      return
    }

    // If the read timed out, send a filler log to keep the log stream alive
    // and get a new connection from the Redis pool.
    if reply == nil {
      readFromTs = redis.LatestXStartTs
      ls.Channel <- StreamKeepAlive
      conn = app.Redis.Get()
      continue
    }

    // Unmarshal raw reply into a BuildableLog.
    log := UnmarshalLog(reply)

    // Stop streaming if completed or failed.
    done = log.Completed || log.Failed

    // Set the next timestamp to read from.
    readFromTs = log.Timestamp

    // Send the latest log through.
    ls.Channel <- log
  }
}