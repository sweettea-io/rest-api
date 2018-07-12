package routes

import (
  "net/http"
  "github.com/gorilla/mux"
  "github.com/sweettea-io/rest-api/app/api/e"
  "github.com/sweettea-io/rest-api/app/api/resp"
  "github.com/sweettea-io/rest-api/defs"
  "github.com/sweettea-io/rest-api/pkg/models"
  "github.com/sweettea-io/rest-api/pkg/utils"
  "github.com/sweettea-io/rest-api/app/api/pl"
  "github.com/sweettea-io/rest-api/app"
  "github.com/lib/pq"
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
  Create a new User (INTERNAL)

  Method:  POST
  Route:   /users
  Payload:
    executor_email    string (required unless using user-creation password for 'executor_password' param)
    executor_password string (required)
    new_email         string (required)
    new_password      string (required)
    admin             bool   (optional, default:false)
*/
func CreateUserHandler(w http.ResponseWriter, req *http.Request) {
  // Validate internal token.
  internalToken := req.Header.Get(defs.AuthHeaderName)

  if internalToken != app.Config.RestApiToken {
    respError(w, e.Unauthorized())
    return
  }

  // Parse & validate payload.
  var payload pl.CreateUserPayload

  // Check if the executor is using the USER_CREATION_HASH to create this user.
  usingUserCreationPw := payload.ExecutorEmail == "" && app.Config.UserCreationHash != "" &&
    utils.VerifySha256(payload.ExecutorPassword, app.Config.UserCreationHash)

  // executor_email can only be empty when using USER_CREATION_HASH as auth method.
  if !payload.Validate(req) || (payload.ExecutorEmail == "" && !usingUserCreationPw) {
    respError(w, e.InvalidPayload())
    return
  }

  // If not using USER_CREATION_HASH for auth, verify executor exists using email/pw.
  if !usingUserCreationPw {
    // Get executor user by email.
    var executorUser models.User
    result := db.Where(&models.User{Email: payload.ExecutorEmail, IsDestroyed: false}).Find(&executorUser)

    if result.RecordNotFound() {
      respError(w, e.UserNotFound())
      return
    }

    // Ensure executor user's password is correct.
    if !utils.VerifyBcrypt(payload.ExecutorPassword, executorUser.HashedPw) {
      respError(w, e.Unauthorized())
      return
    }

    // Only admin users can create other users.
    if !executorUser.Admin {
      logger.Errorf("Error creating new user: executor user %s is not an admin.\n", executorUser.Email)
      respError(w, e.Unauthorized())
      return
    }
  }

  // Hash provided user password.
  hashedPw, err := utils.BcryptHash(payload.NewPassword)

  if err != nil {
    logger.Errorf("Error hashing password during new user creation: %s\n", err.Error())
    respError(w, e.ISE())
    return
  }

  // Create new User.
  newUser := models.User{
    Email: payload.NewEmail,
    HashedPw: hashedPw,
    Admin: payload.Admin,
  }

  if err := db.Create(&newUser).Error; err != nil {
    logger.Errorf("Error creating new user: %s\n", err.Error())

    if err.(*pq.Error).Code.Name() == "unique_violation" {
     respError(w, e.EmailNotAvailable())
    } else {
     respError(w, e.UserCreationFailed())
    }

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
  if !utils.VerifyBcrypt(payload.Password, user.HashedPw) {
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