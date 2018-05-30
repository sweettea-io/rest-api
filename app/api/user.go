package api

import (
  "encoding/json"
  "net/http"
  "github.com/gorilla/mux"
  "github.com/sweettea/rest-api/app/api/e"
  "github.com/sweettea/rest-api/app/api/resp"
  "github.com/sweettea/rest-api/defs"
  "github.com/sweettea/rest-api/pkg/models"
  "github.com/sweettea/rest-api/pkg/utils"
)

// ----------- ROUTER SETUP ------------

const UserRoute = "/users"

func InitUserRouter(baseRouter *mux.Router) {
  // Create user router.
  userRouter := baseRouter.PathPrefix(UserRoute).Subrouter()

  // Attach route handlers.
  userRouter.HandleFunc("/auth", UserAuthHandler).Methods("POST")
}

// ----------- PAYLOADS -----------------

type UserAuthPayload struct {
  Email    string `json:"email"`
  Password string `json:"password"`
}

// ----------- ROUTE HANDLERS -----------

// POST /users/auth
// Basic auth user login
func UserAuthHandler(w http.ResponseWriter, req *http.Request) {
  var payload UserAuthPayload

  // Parse payload and fail if invalid.
  if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
    respError(w, e.InvalidPayload())
  }

  // Get user by email.
  var user models.User
  db.Where(&models.User{Email: payload.Email}).First(&user)

  // Ensure passwords match and fail if not.
  if !utils.VerifyPw(payload.Password, user.HashedPw) {
    respError(w, e.Unauthorized())
    return
  }

  // TODO: Do this in a transaction.
  // Create a new session for the user.
  session := models.Session{User: user}
  db.Create(&session)

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