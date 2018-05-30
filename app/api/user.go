package api

import (
  "net/http"
  "github.com/gorilla/mux"
  "github.com/sweettea/rest-api/app/api/resp"
  "github.com/sweettea/rest-api/defs"
  "github.com/sweettea/rest-api/pkg/models"
  "github.com/sweettea/rest-api/app/api/err"
  "github.com/sweettea/rest-api/pkg/utils"
)

const UserRoute = "/users"

func InitUserRouter(baseRouter *mux.Router) {
  // Create user router.
  userRouter := baseRouter.PathPrefix(UserRoute).Subrouter()

  // Attach route handlers.
  userRouter.HandleFunc("/auth", UserAuthHandler).Methods("GET")
}

func UserAuthHandler(w http.ResponseWriter, req *http.Request) {
  // Get user by email.
  var user models.User
  db.Where(&models.User{Email: "blah"}).First(&user)

  // Ensure passwords match.
  if !utils.VerifyPw("pw", user.HashedPw) {
    respError(w, err.Unauthorized())
    return
  }

  // TODO: Do this in a transaction.
  // Create a new session for the user.
  session := models.Session{User: user}
  db.Create(&session)

  // Put newly minted session's token inside auth header.
  headers := map[string]string{
    defs.AuthHeaderName: session.Token,
  }

  // Respond with success and new header token.
  respOkWithHeaders(w, resp.UserLoginSuccess, headers)
}

func ExampleOfGettingCurrentUser(w http.ResponseWriter, req *http.Request) {
  var user models.User

  // Get current user from session.
  if e := LoadCurrentUser(w, req, &user); e != nil {
    respError(w, err.Unauthorized())
    return
  }
}