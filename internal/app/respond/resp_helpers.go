package respond

import (
  "encoding/json"
  "io"
  "net/http"
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/app/errmsg"
  "github.com/sweettea-io/rest-api/internal/pkg/util/enc"
  "fmt"
)

func Send(w http.ResponseWriter, status int, body string, headers map[string]string) {
  // Write status header.
  w.WriteHeader(status)

  // Write other provided headers.
  for k, v := range headers {
    w.Header().Set(k, v)
  }

  // Write string
  io.WriteString(w, body)
}

func SendJSON(w http.ResponseWriter, status int, data *enc.JSON, headers map[string]string) {
  // Write status header.
  w.WriteHeader(status)

  // Write other provided headers.
  for k, v := range headers {
    w.Header().Set(k, v)
  }

  // Content-Type always set to JSON here -- will override any provided arg headers.
  w.Header().Set("Content-Type", "application/json")

  // Encode and send JSON string.
  if err := json.NewEncoder(w).Encode(data); err != nil {
    io.WriteString(w, errmsg.JsonEncodingError)
  }
}

func Ok(w http.ResponseWriter, data enc.JSON, headers ...map[string]string) {
  data["ok"] = true

  // Pop out extra headers if they exist.
  var extraHeaders map[string]string
  if len(headers) > 0 {
    extraHeaders = headers[0]
  }

  // Send JSON response.
  SendJSON(w, http.StatusOK, &data, extraHeaders)
}

func Created(w http.ResponseWriter, data enc.JSON, headers ...map[string]string) {
  data["ok"] = true

  // Pop out extra headers if they exist.
  var extraHeaders map[string]string
  if len(headers) > 0 {
    extraHeaders = headers[0]
  }

  // Send JSON response.
  SendJSON(w, http.StatusCreated, &data, extraHeaders)
}

func Error(w http.ResponseWriter, err *errmsg.Error, headers ...map[string]string) {
  // Pop out extra headers if they exist.
  var extraHeaders map[string]string
  if len(headers) > 0 {
    extraHeaders = headers[0]
  }

  // Send JSON response.
  SendJSON(w, err.Status, &err.Data, extraHeaders)

  app.Log.Errorf("Request failed with status:%v code:%v message:%q \n",
    err.Status,
    err.Data["code"],
    err.Data["error"],
  )
}

func StreamResources(w http.ResponseWriter) (func(msg string), <-chan bool, bool) {
  // Get response flusher.
  flusher, ok := w.(http.Flusher)
  if !ok {
    return nil, nil, false
  }

  // Get client connection close notifier.
  cn, ok := w.(http.CloseNotifier)
  if !ok {
    return nil, nil, false
  }

  streamLog := func(msg string) {
    fmt.Fprintln(w, msg)
    flusher.Flush()
  }

  return streamLog, cn.CloseNotify(), true
}