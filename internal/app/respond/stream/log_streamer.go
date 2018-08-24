package stream

import (
  "github.com/sweettea-io/rest-api/internal/pkg/redis"
  "github.com/sweettea-io/rest-api/internal/app"
  "net/http"
  "fmt"
)

type Log struct {
  Timestamp string
  Msg       string
  Completed bool
  Failed    bool
}

type logStreamer struct {
  AbstractStreamer
  streamKey        string
  failedLogHandler *func()
}

func NewLogStreamer(w http.ResponseWriter, streamKey string, onFailedLog *func()) (*logStreamer, error) {
  // Initialize new log streamer.
  ls := &logStreamer{
    streamKey: streamKey,
    failedLogHandler: onFailedLog,
  }

  // Initialize abstract streamer.
  if err := ls.initAbstractStreamer(w, map[string]string{
    "Content-Type": "text/plain",   // only working with text logs
    "Transfer-Encoding": "chunked", // we're gonna stream the logs
    "X-Accel-Buffering": "no",      // prevent logs from getting backed up inside Nginx
  }); err != nil {
    return nil, err
  }

  return ls, nil
}

func (ls *logStreamer) Stream() {
  go ls.watchLogs()
  ls.streamToCompletion()
}

func (ls *logStreamer) watchLogs() {
  // Get new connection from Redis pool.
  conn := app.Redis.Get()
  defer conn.Close()

  readFromTs := redis.DefaultXStartTs // start reading from beginning of Redis stream.
  for {
    entries, err := redis.XRead(&conn, ls.streamKey, readFromTs, redis.DefaultReadTimeout)

    // Pipe any unexpected errors and return (if those occur).
    if err != nil {
      app.Log.Errorf("unexpected error while reading log stream: %s\n", err.Error())
      ls.streamLog(&Log{Msg: "unexpected log stream error"})
      return
    }

    // If the read timed out, send a filler log to keep the log stream alive & get fresh Redis connection.
    if entries == nil || len(entries) == 0 {
      ls.streamLog(&Log{Msg: "..."})
      conn = app.Redis.Get()
      continue
    }

    // Iterate over entries since last timestamp.
    for _, entry := range entries {
      // Unmarshal entry into Log structure
      log := unmarshalLog(&entry)

      // Set the next timestamp to read from.
      readFromTs = log.Timestamp

      // Stream latest log.
      ls.streamLog(log)

      // Call log fail handler if log failed.
      if log.Failed && ls.failedLogHandler != nil {
        failHandler := *ls.failedLogHandler
        failHandler()
      }

      // Stop watching logs if failed or completed.
      if log.Failed || log.Completed {
        ls.completeNotifyCh <- true
        return
      }
    }
  }
}

func (ls *logStreamer) streamLog(log *Log) {
  fmt.Fprintln(ls.writer, log.Msg)
  ls.flusher.Flush()
}

func unmarshalLog(entry *redis.XReadEntry) *Log {
  return &Log{
    Timestamp: entry.Timestamp,
    Msg: entry.Args["msg"],
    Completed: entry.Args["complete"] == "true",
    Failed: entry.Args["level"] == "error",
  }
}