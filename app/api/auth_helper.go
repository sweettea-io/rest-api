package api

import (
  "net/http"
  "github.com/sweettea/rest-api/defs"
  "github.com/sweettea/rest-api/app/api/error"
)

func CurrentUser(w http.ResponseWriter, req *http.Request) {
  // Get session token from header.
  if token := req.Header.Get(defs.AuthHeaderName); token == "" {
    respError(w, error.Unauthorized())
  }

  // TODO: Find session by token and user through session

}