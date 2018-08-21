package cloud

import "github.com/sweettea-io/rest-api/internal/app"

var CurrentCloud Cloud

func InitCloud() {
  var err error

  switch app.Config.CloudProvider {
  case AWS:
    CurrentCloud, err = NewAWSCloud()
  }

  if err != nil {
    panic(err)
  }
}