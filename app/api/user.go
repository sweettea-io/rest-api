package api

import (
  "net/http"
  "github.com/gorilla/mux"
  "github.com/sweettea-io/rest-api/app/api/e"
  "github.com/sweettea-io/rest-api/app/api/resp"
  "github.com/sweettea-io/rest-api/defs"
  "github.com/sweettea-io/rest-api/pkg/models"
  "github.com/sweettea-io/rest-api/pkg/utils"
  "github.com/sweettea-io/rest-api/app/api/pl"
)

// ----------- ROUTER SETUP ------------

const UserRoute = "/users" // prefix for all routes in this file

func InitUserRouter(baseRouter *mux.Router) {
  // Create user router.
  userRouter := baseRouter.PathPrefix(UserRoute).Subrouter()

  // Attach route handlers.
  userRouter.HandleFunc("", CreateUserHandler).Methods("POST")
  userRouter.HandleFunc("/auth", UserAuthHandler).Methods("POST")
}

// ----------- ROUTE HANDLERS -----------

/*
  Create a User.

  Method:  POST
  Route:   /users
  Payload:
    email    string (required)
    password string (required)
    admin    bool   (optional, default:false)
*/
func CreateUserHandler(w http.ResponseWriter, req *http.Request) {
  var currentUser models.User

  // Get current user from session.
  if err := LoadCurrentUser(w, req, &currentUser); err != nil {
    logger.Errorln("Error loading current user")
    respError(w, e.Unauthorized())
    return
  }

  // Only admin users can create other users.
  if !currentUser.Admin {
    logger.Errorf("Error creating new user -- current user %s is not an admin.", currentUser.Email)
    respError(w, e.Unauthorized())
    return
  }

  // Parse & validate payload.
  var payload pl.CreateUserPayload

  if !payload.Validate(req) {
    respError(w, e.InvalidPayload())
    return
  }

  // Check availability of this email.
  var emailCount int
  db.Where(&models.User{Email: payload.Email, IsDestroyed: false}).Count(&emailCount)

  if emailCount > 0 {
    respError(w, e.EmailNotAvailable())
    return
  }

  // Hash provided user password.
  hashedPw, err := utils.HashPw(payload.Password)

  if err != nil {
    logger.Errorf("Error hashing password during new user creation: %s\n", err.Error())
    respError(w, e.ISE())
    return
  }

  // Create new User.
  newUser := models.User{
    Email: payload.Email,
    HashedPw: hashedPw,
    Admin: payload.Admin,
  }

  if err := db.Create(&newUser).Error; err != nil {
    logger.Errorf("Error creating new user: %s\n", err.Error())
    respError(w, e.UserCreationFailed())
    return
  }

  // Create response payload and respond.
  respData := resp.UserCreationSuccess
  respData["uid"] = newUser.Uid

  respOk(w, respData)
}

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
    respError(w, e.SessionCreationFailed())
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