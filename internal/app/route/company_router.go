package route

import (
  "net/http"
  "github.com/lib/pq"
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/app/errmsg"
  "github.com/sweettea-io/rest-api/internal/app/payload"
  "github.com/sweettea-io/rest-api/internal/app/respond"
  "github.com/sweettea-io/rest-api/internal/app/successmsg"
  "github.com/sweettea-io/rest-api/internal/pkg/model"
  "github.com/sweettea-io/rest-api/internal/pkg/util/crypt"
)

// ----------- ROUTER SETUP ------------

const CompanyRoute = "/companies" // prefix for all routes in this file

func InitCompanyRouter() {
  // Create company router.
  companyRouter := Router.PathPrefix(CompanyRoute).Subrouter()

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
  internalToken := req.Header.Get(app.Config.AuthHeaderName)

  if internalToken != app.Config.RestApiToken {
    respond.Error(w, errmsg.Unauthorized())
    return
  }

  // Parse & validate payload.
  var pl payload.CreateCompanyPayload

  if !pl.Validate(req) {
    respond.Error(w, errmsg.InvalidPayload())
  }

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

  // Only admin users can create companies.
  if !executorUser.Admin {
    app.Log.Errorf("Error creating new company: executor user %s is not an admin.\n", executorUser.Email)
    respond.Error(w, errmsg.Unauthorized())
    return
  }

  // Create Company.
  company := model.Company{Name: pl.Name}

  if err := app.DB.Create(&company).Error; err != nil {
    app.Log.Errorf("Error creating new company: %s\n", err.Error())

    if err.(*pq.Error).Code.Name() == "unique_violation" {
      respond.Error(w, errmsg.CompanyAlreadyExists())
    } else {
      respond.Error(w, errmsg.CompanyCreationFailed())
    }

    return
  }

  // Create response pl and respond.
  respData := successmsg.CompanyCreationSuccess
  respData["uid"] = company.Uid

  respond.Created(w, respData)
}