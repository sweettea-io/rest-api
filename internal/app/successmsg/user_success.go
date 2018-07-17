package successmsg

import "github.com/sweettea-io/rest-api/internal/pkg/util/enc"

var UserCreationSuccess = enc.JSON{"message": "User Creation Successful"}
var UserLoginSuccess = enc.JSON{"message": "Login Successful"}