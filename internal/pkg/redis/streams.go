package redis

import (
  r "github.com/gomodule/redigo/redis"
  "net"
  "time"
  "fmt"
  "github.com/sweettea-io/rest-api/internal/pkg/util/maputil"
)

const DefaultReadTimeout = 30 // seconds
const DefaultXStartTs = "0.0" // very beginning of stream

type XReadEntry struct {
  Timestamp string
  Args      map[string]string
}

func XRead(conn *r.Conn, stream string, startTs string, timeout int) ([]XReadEntry, error) {
  reply, err := r.Values(r.DoWithTimeout(
    *conn,
    time.Duration(timeout) * time.Second,
    "XREAD",
    "BLOCK",
    timeout * 1000,
    "STREAMS",
    stream,
    startTs,
  ))

  // If no error, return parsed reply.
  if err == nil {
    parsedReply, err := parseXReadReply(reply)

    if err != nil {
      return nil, err
    }

    return parsedReply[stream], nil
  }

  // If timeout error, return nil for both.
  if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
    return nil, nil
  }

  // Lastly, an unexpected error occurred (not a timeout error), so return that.
  return nil, err
}

func parseXReadReply(streams []interface{}) (map[string][]XReadEntry, error) {
  xReadResp := map[string][]XReadEntry{}

  for _, s := range streams {
    stream, ok := s.([]interface{})
    if !ok {
      return nil, fmt.Errorf("error converting interface{} to []interface{} during XREAD parse")
    }

    // Parse stream name and entries.
    var streamName string
    var streamEntries []interface{}
    if _, err := r.Scan(stream, &streamName, &streamEntries); err != nil {
      return nil, err
    }

    // Initialize streamName as key in response object.
    xReadResp[streamName] = []XReadEntry{}

    // Add entries to response map.
    for _, se := range streamEntries {
      streamEntry, ok := se.([]interface{})
      if !ok {
        return nil, fmt.Errorf("error converting interface{} to []interface{} during XREAD parse")
      }

      var timestamp string
      var args []string
      if _, err := r.Scan(streamEntry, &timestamp, &args); err != nil {
        return nil, err
      }

      // Create new XReadEntry.
      xReadEntry := XReadEntry{
        Timestamp: timestamp,
        Args: maputil.FromStrSlice(args),
      }

      // Append entry to stream key.
      xReadResp[streamName] = append(xReadResp[streamName], xReadEntry)
    }
  }

  return xReadResp, nil
}