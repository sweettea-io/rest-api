package stream

import (
  "net/http"
  "fmt"
)

type AbstractStreamer struct {
  writer           http.ResponseWriter
  flusher          http.Flusher
  closeNotifyCh    <-chan bool
  completeNotifyCh <-chan bool
  respHeaders      map[string]string
}

func (as *AbstractStreamer) initAbstractStreamer(w http.ResponseWriter, respHeaders map[string]string) error {
  // Get response flusher from writer.
  flusher, ok := w.(http.Flusher)
  if !ok {
    return fmt.Errorf("flusher not created -- streaming not supported")
  }

  // Get client connection close notifier.
  closeNotifier, ok := w.(http.CloseNotifier)
  if !ok {
    return fmt.Errorf("closeNotifier not created -- streaming not supported")
  }

  // Set abstract fields.
  as.writer = w
  as.flusher = flusher
  as.closeNotifyCh = closeNotifier.CloseNotify()
  as.completeNotifyCh = make(chan bool)
  as.respHeaders = respHeaders

  return nil
}

func (as *AbstractStreamer) applyRespHeaders() {
  for k, v := range as.respHeaders {
    as.writer.Header().Set(k, v)
  }
}

func (as *AbstractStreamer) streamToCompletion() {
  // Only apply response headers once we're ready to begin streaming.
  as.applyRespHeaders()

  for {
    select {
    // Return if client closes the connection.
    case <-as.closeNotifyCh:
      return

      // Return when generator is complete.
    case <-as.completeNotifyCh:
      return
    }
  }
}