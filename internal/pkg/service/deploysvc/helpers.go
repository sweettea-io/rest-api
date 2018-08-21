package deploysvc

import (
  "fmt"
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/pkg/model"
)

func NewHostname(deploy *model.Deploy) string {
  return fmt.Sprintf("%s.%s", deploy.Slug, app.Config.Domain)
}