package route

import (
  "net/http"
  "github.com/lib/pq"
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/app/errmsg"
  "github.com/sweettea-io/rest-api/internal/app/payload"
  "github.com/sweettea-io/rest-api/internal/app/respond"
  "github.com/sweettea-io/rest-api/internal/app/successmsg"
  "github.com/sweettea-io/rest-api/internal/pkg/service/sessionsvc"
  "github.com/sweettea-io/rest-api/internal/pkg/service/usersvc"
  "github.com/sweettea-io/rest-api/internal/pkg/util/crypt"
)

// ----------- ROUTER SETUP ------------

// Prefix for all routes in this file.
const UserRoute = "/user"

func InitUserRouter() {
  // Create user router.
  userRouter := Router.GetRouter().PathPrefix(UserRoute).Subrouter()

  // Attach route handlers.
  userRouter.HandleFunc("", CreateUserHandler).Methods("POST")
  userRouter.HandleFunc("/auth", UserAuthHandler).Methods("POST")
}

// ----------- ROUTE HANDLERS -----------

/*
  Create a new User (INTERNAL)

  Method:  POST
  Route:   /user
  Payload:
    executorEmail     string (required unless using user-creation password for 'executorPassword' param)
    executorPassword  string (required)
    newEmail          string (required)
    newPassword       string (required)
    admin             bool   (optional, default:false)
*/
func CreateUserHandler(w http.ResponseWriter, req *http.Request) {
  // Validate internal token.
  if internalToken := req.Header.Get(app.Config.AuthHeaderName); internalToken != app.Config.RestApiToken {
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
    executorUser, err := usersvc.FromEmail(pl.ExecutorEmail)

    if err != nil {
      app.Log.Errorln(err.Error())
      respond.Error(w, errmsg.UserNotFound())
      return
    }

    // Ensure executor user's password is correct.
    if !crypt.VerifyBcrypt(pl.ExecutorPassword, executorUser.HashedPw) {
      app.Log.Errorln("error creating new User: invalid executor user password")
      respond.Error(w, errmsg.Unauthorized())
      return
    }

    // Only admin users can create other users.
    if !executorUser.Admin {
      app.Log.Errorln("error creating new User: executor user must be an admin")
      respond.Error(w, errmsg.Unauthorized())
      return
    }
  }

  // Hash provided user password.
  hashedPw, err := crypt.BcryptHash(pl.NewPassword)

  if err != nil {
    app.Log.Errorf("error creating new User: bcrypt password hash failed with %s\n", err.Error())
    respond.Error(w, errmsg.ISE())
    return
  }

  // Create new User.
  newUser, err := usersvc.Create(pl.NewEmail, hashedPw, pl.Admin)

  if err != nil {
    app.Log.Errorln(err.Error())
    pqError, ok := err.(*pq.Error)

    if ok && pqError.Code.Name() == "unique_violation" {
      respond.Error(w, errmsg.EmailNotAvailable())
    } else {
      respond.Error(w, errmsg.UserCreationFailed())
    }

    return
  }

  // Create response payload and respond.
  respData := successmsg.UserCreationSuccess
  respData["uid"] = newUser.Uid

  respond.Created(w, respData)
}

/*
  User login with basic auth.
  If successful, returns newly minted session token inside header.

  Method:  POST
  Route:   /user/auth
  Payload:
    email     string (required)
    password  string (required)
*/
func UserAuthHandler(w http.ResponseWriter, req *http.Request) {
  // Parse & validate payload.
  var pl payload.UserAuthPayload

  if !pl.Validate(req) {
    respond.Error(w, errmsg.InvalidPayload())
    return
  }

  // Get user by email.
  user, err := usersvc.FromEmail(pl.Email)

  if err != nil {
    app.Log.Errorln(err.Error())
    respond.Error(w, errmsg.UserNotFound())
    return
  }

  // Ensure passwords match.
  if !crypt.VerifyBcrypt(pl.Password, user.HashedPw) {
    app.Log.Errorln("user auth error: invalid password")
    respond.Error(w, errmsg.Unauthorized())
    return
  }

  // Create new session for user.
  session, err := sessionsvc.Create(user)

  if err != nil {
    app.Log.Errorln(err.Error())
    respond.Error(w, errmsg.SessionCreationFailed())
    return
  }

  // Put newly minted session token inside auth header.
  headers := map[string]string{
    app.Config.AuthHeaderName: session.Token,
  }

  // Respond with success and new auth token.
  respond.Ok(w, successmsg.UserLoginSuccess, headers)
}