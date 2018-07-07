package routes

import (
  "encoding/json"
  "io"
  "net/http"
  "github.com/gorilla/mux"
  "github.com/jinzhu/gorm"
  "github.com/Sirupsen/logrus"
  "github.com/sweettea-io/rest-api/app/api/e"
  "github.com/sweettea-io/rest-api/app/api/mw"
  "github.com/sweettea-io/rest-api/app/api/helpers"
  "github.com/sweettea-io/rest-api/pkg/utils"
)

// Create global vars for our db and logger so that all other
// routes in this package can reference them.
var db *gorm.DB
var logger *logrus.Logger

// Create base router and attach all subroutes.
func CreateRouter(baseRoute string, database *gorm.DB, l *logrus.Logger) *mux.Router {
  // Assign values to our global vars declared above.
  db = database
  logger = l

  // Init middleware and helpers with global vars.
  mw.Init(db, logger)
  helpers.Init(db, logger)

  // Create base router from provided baseRoute.
  baseRouter := mux.NewRouter().PathPrefix(baseRoute).Subrouter()

  // Attach base middleware.
  baseRouter.Use(mw.LogRequest)

  // Create route groups for each model needing a REST-ful interface.
  InitUserRouter(baseRouter)
  InitCompanyRouter(baseRouter)
  InitClusterRouter(baseRouter)

  return baseRouter
}

// ---------- Route-agnostic utilities ----------

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