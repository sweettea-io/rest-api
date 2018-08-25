package logger

import (
  "github.com/Sirupsen/logrus"
  "fmt"
  r "github.com/gomodule/redigo/redis"
)

// Define different log levels for this logger:
var (
  // InfoLevel applies to general operational entries about what's going on inside the app.
  InfoLevel = logrus.InfoLevel.String()

  // WarnLevel applies to non-critical entries that deserve eyes.
  WarnLevel = logrus.WarnLevel.String()

  // ErrorLevel applies to errors that should definitely be noted.
  ErrorLevel = logrus.ErrorLevel.String()
)

// Lgr is a wrapper type around `*logrus.Logger`, providing Redis stream functionality.
type Lgr struct {
  Logger *logrus.Logger
  RedisPool *r.Pool
  Stream string
}

// TODO: Update the comments here

// Info logs to the InfoLevel.
func (l *Lgr) Info(args ...interface{}) {
  l.Logger.Info(args...)
}

// Infof logs a formatted string to the InfoLevel.
func (l *Lgr) Infof(format string, args ...interface{}) {
  l.Logger.Infof(format, args...)
}

// Infoln is equivalent to `Info` but appends a new line to the message.
func (l *Lgr) Infoln(args ...interface{}) {
  l.Logger.Infoln(args...)
}

// BuildableInfo logs to the InfoLevel and the Redis stream of the buildable.
func (l *Lgr) BuildableInfo(buildableUid string, args ...interface{}) {
  msg := fmt.Sprint(args...)
  l.Logger.Info(msg)
  l.newStreamEntry(buildableUid, msg, InfoLevel, false)
}

// BuildableInfof logs a formatted string to the InfoLevel and the Redis stream of the buildable.
func (l *Lgr) BuildableInfof(buildableUid string, format string, args ...interface{}) {
  msg := fmt.Sprintf(format, args...)
  l.Logger.Info(msg)
  l.newStreamEntry(buildableUid, msg, InfoLevel, false)
}

// BuildableInfoln is equivalent to `BuildableInfo` but appends a new line to the message.
func (l *Lgr) BuildableInfoln(buildableUid string, args ...interface{}) {
  msg := fmt.Sprintln(args...)
  l.Logger.Info(msg)
  l.newStreamEntry(buildableUid, msg, InfoLevel, false)
}

func (l *Lgr) CompleteBuildable(buildableUid string, args ...interface{}) {
  msg := fmt.Sprint(args...)
  l.Logger.Info(msg)
  l.newStreamEntry(buildableUid, msg, InfoLevel, true)
}

func (l *Lgr) CompleteBuildablef(buildableUid string, format string, args ...interface{}) {
  msg := fmt.Sprintf(format, args...)
  l.Logger.Info(msg)
  l.newStreamEntry(buildableUid, msg, InfoLevel, true)
}

func (l *Lgr) CompleteBuildableln(buildableUid string, args ...interface{}) {
  msg := fmt.Sprintln(args...)
  l.Logger.Info(msg)
  l.newStreamEntry(buildableUid, msg, InfoLevel, true)
}

// Warn logs to the WarnLevel.
func (l *Lgr) Warn(args ...interface{}) {
  l.Logger.Warn(args...)
}

// Warnf logs a formatted string to the WarnLevel.
func (l *Lgr) Warnf(format string, args ...interface{}) {
  l.Logger.Warnf(format, args...)
}

// Warnln is equivalent to `Warn` but appends a new line to the message.
func (l *Lgr) Warnln(args ...interface{}) {
  l.Logger.Warnln(args...)
}

// BuildableWarn logs to the WarnLevel and the Redis stream of the buildable.
func (l *Lgr) BuildableWarn(buildableUid string, args ...interface{}) {
  msg := fmt.Sprint(args...)
  l.Logger.Warn(msg)
  l.newStreamEntry(buildableUid, msg, WarnLevel, false)
}

// BuildableWarnf logs a formatted string to the WarnLevel and the Redis stream of the buildable.
func (l *Lgr) BuildableWarnf(buildableUid string, format string, args ...interface{}) {
  msg := fmt.Sprintf(format, args...)
  l.Logger.Warn(msg)
  l.newStreamEntry(buildableUid, msg, WarnLevel, false)
}

// BuildableWarnln is equivalent to `BuildableWarn` but appends a new line to the message.
func (l *Lgr) BuildableWarnln(buildableUid string, args ...interface{}) {
  msg := fmt.Sprintln(args...)
  l.Logger.Warn(msg)
  l.newStreamEntry(buildableUid, msg, WarnLevel, false)
}

// Error logs to the ErrorLevel.
func (l *Lgr) Error(args ...interface{}) {
  l.Logger.Error(args...)
}

// Errorf logs a formatted string to the ErrorLevel.
func (l *Lgr) Errorf(format string, args ...interface{}) {
  l.Logger.Errorf(format, args...)
}

// Errorln is equivalent to `Error` but appends a new line to the message.
func (l *Lgr) Errorln(args ...interface{}) {
  l.Logger.Errorln(args...)
}

// BuildableError logs to the ErrorLevel and the Redis stream of the buildable.
func (l *Lgr) BuildableError(buildableUid string, args ...interface{}) {
  msg := fmt.Sprint(args...)
  l.Logger.Error(msg)
  l.newStreamEntry(buildableUid, msg, ErrorLevel, false)
}

// BuildableErrorf logs a formatted string to the ErrorLevel and the Redis stream of the buildable.
func (l *Lgr) BuildableErrorf(buildableUid string, format string, args ...interface{}) {
  msg := fmt.Sprintf(format, args...)
  l.Logger.Error(msg)
  l.newStreamEntry(buildableUid, msg, ErrorLevel, false)
}

// BuildableErrorln is equivalent to `BuildableError` but appends a new line to the message.
func (l *Lgr) BuildableErrorln(buildableUid string, args ...interface{}) {
  msg := fmt.Sprintln(args...)
  l.Logger.Error(msg)
  l.newStreamEntry(buildableUid, msg, ErrorLevel, false)
}

// newStreamEntry adds a new message to the logger's redis stream.
func (l *Lgr) newStreamEntry(streamKey string, msg string, level string, complete bool) {
  // Get new connection from Redis pool.
  conn := l.RedisPool.Get()
  defer conn.Close()

  // Add message to log stream.
  if _, err := conn.Do("XADD", streamKey, "*", "msg", msg, "level", level, "complete", complete); err != nil {
    l.Errorf("error logging to stream: %s", err.Error())
  }
}

func NewLogger(pool *r.Pool) *Lgr {
  return &Lgr{
    Logger: logrus.New(),
    RedisPool: pool,
  }
}