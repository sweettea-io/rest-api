package stream

import (
  "github.com/sweettea-io/rest-api/internal/pkg/redis"
  "github.com/sweettea-io/rest-api/internal/app"
  "net/http"
  "fmt"
  "github.com/sweettea-io/rest-api/internal/pkg/util/typeconvert"
  "github.com/sweettea-io/rest-api/internal/pkg/logger"
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

// Log message sent to keep stream alive.
var StreamKeepAlive = &Log{Msg: "..."}

// Log message to stream to user when unexpected error occurs.
var StreamErrorLog = &Log{Msg: "unexpected log stream error"}

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

    // Pipe any unexpected errors and return (if one occurred).
    if err != nil {
      ls.handleUnexpectedErr("unexpected error while reading log stream: %s\n", err.Error())
      return
    }

    // If the read timed out, send a filler log to keep the log stream alive & get fresh Redis connection.
    if entries == nil || len(entries) == 0 {
      ls.streamLog(StreamKeepAlive)
      conn = app.Redis.Get()
      continue
    }

    // Iterate over entries since last timestamp.
    for _, entry := range entries {
      // Unmarshal entry into Log structure.
      log, err := unmarshalLog(&entry)

      // Pipe any unexpected errors and return (if one occurred).
      if err != nil {
        ls.handleUnexpectedErr("error unmarshalling redis stream reply into Log: %s\n", err.Error())
        return
      }

      // Set the next timestamp to read from.
      readFromTs = log.Timestamp

      // Stream latest log.
      ls.streamLog(log)

      // Call log fail handler if log failed.
      if log.Failed {
        ls.handleFailure()
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

func (ls *logStreamer) handleFailure() {
  if ls.failedLogHandler == nil {
    return
  }

  failHandler := *ls.failedLogHandler
  failHandler()
}

func (ls *logStreamer) handleUnexpectedErr(format string, args ...interface{}) {
  app.Log.Errorf(format, args...)
  ls.streamLog(StreamErrorLog)
  ls.handleFailure()
  ls.completeNotifyCh <- true
}

func unmarshalLog(entry *redis.XReadEntry) (*Log, error) {
  // Parse log message from args.
  msgBytes, err := typeconvert.InterfaceToBytes(entry.Args["msg"])
  if err != nil {
    return nil, err
  }
  msg := typeconvert.BytesToStr(msgBytes)

  // Parse log level from args.
  levelBytes, err := typeconvert.InterfaceToBytes(entry.Args["level"])
  if err != nil {
    return nil, err
  }
  level := typeconvert.BytesToStr(levelBytes)

  // Parse complete status from args.
  complete, err := typeconvert.InterfaceToBool(entry.Args["complete"])
  if err != nil {
    return nil, err
  }

  // Create new Log instance.
  log := &Log{
    Timestamp: entry.Timestamp,
    Msg: msg,
    Completed: complete,
    Failed: level == logger.ErrorLevel,
  }

  return log, nil
}