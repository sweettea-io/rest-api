package resp

import "github.com/sweettea-io/rest-api/pkg/utils"

// ------- User Responses --------

var UserLoginSuccess = utils.JSON{"message": "Login Successful"}
var UserCreationSuccess = utils.JSON{"message": "User Creation Successful"}

// ------- Company Responses --------

var CompanyCreationSuccess = utils.JSON{"message": "Company Creation Successful"}