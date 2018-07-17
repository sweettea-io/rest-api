package route

import (
  "net/http"
  "github.com/lib/pq"
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/app/respond"
  "github.com/sweettea-io/rest-api/internal/app/errmsg"
  "github.com/sweettea-io/rest-api/internal/pkg/util/crypt"
  "github.com/sweettea-io/rest-api/internal/app/payload"
  "github.com/sweettea-io/rest-api/internal/pkg/model"
  "github.com/sweettea-io/rest-api/internal/app/successmsg"
)

// ----------- ROUTER SETUP ------------

// Prefix for all routes in this file
const UserRoute = "/users"

func InitUserRouter() {
  // Create user router.
  userRouter := Router.PathPrefix(UserRoute).Subrouter()

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
  internalToken := req.Header.Get(app.Config.AuthHeaderName)

  if internalToken != app.Config.RestApiToken {
    respond.Error(w, errmsg.Unauthorized())
    return
  }

  // Parse & validate payload.
  var pl payload.CreateUserPayload

  if !pl.Validate(req) {
    respond.Error(w, errmsg.InvalidPayload())
    return
  }

  // Check if the executor is using the USER_CREATION_HASH to create this user.
  usingUserCreationPw := pl.ExecutorEmail == "" && app.Config.UserCreationHash != "" &&
    crypt.VerifySha256(pl.ExecutorPassword, app.Config.UserCreationHash)

  // If not using USER_CREATION_HASH for auth, verify executor exists using email/pw.
  if !usingUserCreationPw {
    // Get executor user by email.
    var executorUser model.User
    result := app.DB.Where(&model.User{Email: pl.ExecutorEmail}).Find(&executorUser)

    if result.RecordNotFound() {
      respond.Error(w, errmsg.UserNotFound())
      return
    }

    // Ensure executor user's password is correct.
    if !crypt.VerifyBcrypt(pl.ExecutorPassword, executorUser.HashedPw) {
      respond.Error(w, errmsg.Unauthorized())
      return
    }

    // Only admin users can create other users.
    if !executorUser.Admin {
      app.Log.Errorf("Error creating new user: executor user %s is not an admin.\n", executorUser.Email)
      respond.Error(w, errmsg.Unauthorized())
      return
    }
  }

  // Hash provided user password.
  hashedPw, err := crypt.BcryptHash(pl.NewPassword)

  if err != nil {
    app.Log.Errorf("Error hashing password during new user creation: %s\n", err.Error())
    respond.Error(w, errmsg.ISE())
    return
  }

  // Create new User.
  newUser := model.User{
    Email: pl.NewEmail,
    HashedPw: hashedPw,
    Admin: pl.Admin,
  }

  if err := app.DB.Create(&newUser).Error; err != nil {
    app.Log.Errorf("Error creating new user: %s\n", err.Error())

    if err.(*pq.Error).Code.Name() == "unique_violation" {
      respond.Error(w, errmsg.EmailNotAvailable())
    } else {
      respond.Error(w, errmsg.UserCreationFailed())
    }

    return
  }

  // Create response payload and respond.
  respData := successmsg.UserCreationSuccess
  respData["uid"] = newUser.Uid

  respond.Ok(w, respData)
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
  var pl payload.UserAuthPayload

  if !pl.Validate(req) {
    respond.Error(w, errmsg.InvalidPayload())
    return
  }

  // Get user by email.
  var user model.User
  result := app.DB.Where(&model.User{Email: pl.Email}).First(&user)

  // Ensure user exists.
  if result.RecordNotFound() {
    respond.Error(w, errmsg.UserNotFound())
    return
  }

  // Ensure passwords match.
  if !crypt.VerifyBcrypt(pl.Password, user.HashedPw) {
    respond.Error(w, errmsg.Unauthorized())
    return
  }

  // Create new session for user.
  session := model.Session{User: user}

  if err := app.DB.Create(&session).Error; err != nil {
    respond.Error(w, errmsg.SessionCreationFailed())
    app.Log.Errorf("Session creation failed for User(id=%v): %s\n", user.ID, err.Error())
    return
  }

  // Put newly minted session token inside auth header.
  headers := map[string]string{
    app.Config.AuthHeaderName: session.Token,
  }

  // Respond with success and new auth token.
  respond.Ok(w, successmsg.UserLoginSuccess, headers)
}