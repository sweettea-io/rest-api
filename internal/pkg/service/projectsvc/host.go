package projectsvc

import (
  "strings"
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/pkg/model"
  "github.com/sweettea-io/rest-api/internal/pkg/util/projecthost"
)

func GetHost(project *model.Project) projecthost.Host {
  var host projecthost.Host

  switch hostNameForNsp(project.Nsp) {
  case projecthost.GitHubName:
    host = &projecthost.GitHub{Token: app.Config.GitHubAccessToken}
  default:
    return nil
  }

  host.Configure()

  return host
}

func hostNameForNsp(nsp string) string {
  switch true {
  case strings.HasPrefix(nsp, projecthost.GitHubDomain):
    return projecthost.GitHubName
  default:
    return ""
  }
}