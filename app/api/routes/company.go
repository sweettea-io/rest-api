package routes

import (
  "github.com/gorilla/mux"
  "net/http"
  "github.com/sweettea-io/rest-api/app"
  "github.com/sweettea-io/rest-api/defs"
  "github.com/sweettea-io/rest-api/app/api/e"
  "github.com/sweettea-io/rest-api/app/api/pl"
  "github.com/sweettea-io/rest-api/pkg/models"
  "github.com/sweettea-io/rest-api/pkg/utils"
  "github.com/lib/pq"
  "github.com/sweettea-io/rest-api/app/api/resp"
)

// ----------- ROUTER SETUP ------------

const CompanyRoute = "/companies" // prefix for all routes in this file

func InitCompanyRouter(baseRouter *mux.Router) {
  // Create company router.
  companyRouter := baseRouter.PathPrefix(CompanyRoute).Subrouter()

  // Attach route handlers.
  companyRouter.HandleFunc("", CreateCompanyHandler).Methods("POST")
}

// ----------- ROUTE HANDLERS -----------

/*
  Create a new Company (INTERNAL)

  Method:  POST
  Route:   /companies
  Payload:
    executor_email    string (required)
    executor_password string (required)
    name              string (required)
*/
func CreateCompanyHandler(w http.ResponseWriter, req *http.Request) {
  // Validate internal token.
  internalToken := req.Header.Get(defs.AuthHeaderName)

  if internalToken != app.Config.RestApiToken {
    respError(w, e.Unauthorized())
    return
  }

  // Parse & validate payload.
  var payload pl.CreateCompanyPayload

  if !payload.Validate(req) {
    respError(w, e.InvalidPayload())
  }

  // Get executor user by email.
  var executorUser models.User
  result := db.Where(&models.User{Email: payload.ExecutorEmail, IsDestroyed: false}).Find(&executorUser)

  if result.RecordNotFound() {
    respError(w, e.UserNotFound())
    return
  }

  // Ensure executor user's password is correct.
  if !utils.VerifyPw(payload.ExecutorPassword, executorUser.HashedPw) {
    respError(w, e.Unauthorized())
    return
  }

  // Only admin users can create companies.
  if !executorUser.Admin {
    logger.Errorf("Error creating new company: executor user %s is not an admin.\n", executorUser.Email)
    respError(w, e.Unauthorized())
    return
  }

  // Create Company.
  company := models.Company{Name: payload.Name}

  if err := db.Create(&company).Error; err != nil {
    logger.Errorf("Error creating new company: %s\n", err.Error())

    if err.(*pq.Error).Code.Name() == "unique_violation" {
      respError(w, e.CompanyAlreadyExists())
    } else {
      respError(w, e.CompanyCreationFailed())
    }

    return
  }

  // Create response payload and respond.
  respData := resp.CompanyCreationSuccess
  respData["uid"] = company.Uid

  respOk(w, respData)
}