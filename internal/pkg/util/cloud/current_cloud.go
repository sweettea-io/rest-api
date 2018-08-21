package cloud

import "github.com/sweettea-io/rest-api/internal/app"

var CurrentCloud Cloud

func InitCloud() {
  switch app.Config.CloudProvider {
  case AWS:
    CurrentCloud = &AWSCloud{}
  default:
    CurrentCloud = nil
  }

  // This is okay for environments where CLOUD_PROVIDER env var isn't required.
  if CurrentCloud == nil {
    app.Log.Warnln("Not initializing current cloud session.")
    return
  }

  // Initialize current cloud.
  CurrentCloud.Init()
}