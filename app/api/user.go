package api

import (
  "net/http"
  "github.com/gorilla/mux"
  "github.com/sweettea/rest-api/app/api/e"
  "github.com/sweettea/rest-api/app/api/resp"
  "github.com/sweettea/rest-api/defs"
  "github.com/sweettea/rest-api/pkg/models"
  "github.com/sweettea/rest-api/pkg/utils"
  "github.com/sweettea/rest-api/app/api/pl"
)

// ----------- ROUTER SETUP ------------

const UserRoute = "/users" // prefix for all routes in this file

func InitUserRouter(baseRouter *mux.Router) {
  // Create user router.
  userRouter := baseRouter.PathPrefix(UserRoute).Subrouter()

  // Attach route handlers.
  userRouter.HandleFunc("/auth", UserAuthHandler).Methods("POST")
}

// ----------- ROUTE HANDLERS -----------

/*
  User login with basic auth.

  Method:  POST
  Route:   /users/auth
  Payload:
    email    string (required)
    password string (required)
 */
func UserAuthHandler(w http.ResponseWriter, req *http.Request) {
  // Parse & validate payload.
  var payload pl.UserAuthPayload

  if !payload.Validate(req) {
    respError(w, e.InvalidPayload())
    return
  }

  // Get user by email.
  var user models.User
  result := db.Where(&models.User{Email: payload.Email, IsDestroyed: false}).First(&user)

  // Ensure user exists.
  if result.RecordNotFound() {
    respError(w, e.UserNotFound())
    return
  }

  // Ensure passwords match.
  if !utils.VerifyPw(payload.Password, user.HashedPw) {
    respError(w, e.Unauthorized())
    return
  }

  // Create new session for user.
  session := models.Session{User: user}

  if err := db.Create(&session).Error; err != nil {
    respError(w, e.ISE())
    logger.Errorf("Session creation failed for User(id=%v): %s\n", user.ID, err.Error())
    return
  }

  // Put newly minted session token inside auth header.
  headers := map[string]string{
    defs.AuthHeaderName: session.Token,
  }

  // Respond with success and new auth token.
  respOkWithHeaders(w, resp.UserLoginSuccess, headers)
}

func ExampleOfGettingCurrentUser(w http.ResponseWriter, req *http.Request) {
  var user models.User

  // Get current user from session.
  if err := LoadCurrentUser(w, req, &user); err != nil {
    respError(w, e.Unauthorized())
    return
  }
}