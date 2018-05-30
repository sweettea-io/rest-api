package api

import (
  "net/http"
  "github.com/sweettea/rest-api/app/api/e"
  "github.com/sweettea/rest-api/pkg/utils"
  "encoding/json"
  "io"
)

func respJson(w http.ResponseWriter, status int, data *utils.JSON) {
  // Write status header.
  w.WriteHeader(status)

  // All about that JSON
  w.Header().Set("Content-Type", "application/json")

  // Encode and send JSON string.
  if err := json.NewEncoder(w).Encode(data); err != nil {
    io.WriteString(w, e.JsonEncodingError)
  }
}

func respCreated(w http.ResponseWriter, data utils.JSON) {
  data["ok"] = true
  respJson(w, http.StatusCreated, &data)
}

func respOk(w http.ResponseWriter, data utils.JSON) {
  data["ok"] = true
  respJson(w, http.StatusOK, &data)
}

func respOkWithHeaders(w http.ResponseWriter, data utils.JSON, headers map[string]string) {
  // Add each provided header from the provided map.
  for k, v := range headers {
    w.Header().Set(k, v)
  }

  // Continue with regular OK response.
  respOk(w, data)
}

func respError(w http.ResponseWriter, err *e.Error) {
  respJson(w, err.Status, &err.Data)

  logger.Errorf("Request failed with status:%v code:%v message:%q \n",
    err.Status,
    err.Data["code"],
    err.Data["error"],
  )
}