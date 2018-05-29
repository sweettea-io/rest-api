package api

import (
  "io"
  "net/http"
  "github.com/sweettea/rest-api/app/api/error"
  "github.com/sweettea/rest-api/pkg/utils"
)

func respJson(w http.ResponseWriter, status int, data *utils.JSON) {
  w.WriteHeader(status)
  w.Header().Set("Content-Type", "application/json")
  io.WriteString(w, data.Stringify())
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

func respError(w http.ResponseWriter, error *error.Error) {
  respJson(w, error.Status, &error.Data)
}