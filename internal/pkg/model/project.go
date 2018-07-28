package model

import (
  "github.com/jinzhu/gorm"
  "github.com/sweettea-io/rest-api/internal/pkg/util/unique"
  "github.com/sweettea-io/rest-api/internal/pkg/util/enc"
  "github.com/sweettea-io/rest-api/internal/pkg/util/projecthost"
  "fmt"
  "strings"
)

/*
  has_one -- ProjectConfig

  has_many --> Commits
  has_many --> Models
*/
type Project struct {
  gorm.Model
  Uid             string        `gorm:"type:varchar(240);default:null;not null;unique;index:project_uid"`
  Host            string        `gorm:"type:varchar(240);default:null;not null"`
  Nsp             string        `gorm:"type:varchar(360);default:null;not null;unique;index:project_nsp"`
  ProjectConfig   *ProjectConfig
  ProjectConfigID uint          `gorm:"default:null;not null;index:project_config_id"`
  Commits         []Commit
  Models          []Model
}

// Assign Host to Project before saving.
func (project *Project) BeforeSave(scope *gorm.Scope) error {
  // Downcase project namespace.
  scope.SetColumn("Nsp", strings.ToLower(project.Nsp))

  // Get supported host for namespace.
  host := hostNameForNsp(project.Nsp)

  if host == "" {
    return fmt.Errorf("invalid project namespace \"%s\" -- no recognizable host", project.Nsp)
  }

  scope.SetColumn("Host", host)
  return nil
}

// Assign Uid to Project before creation.
func (project *Project) BeforeCreate(scope *gorm.Scope) error {
  scope.SetColumn("Uid", unique.NewUid())
  return nil
}

func (project *Project) AsJSON() enc.JSON {
  return enc.JSON{
    "uid": project.Uid,
    "host": project.Host,
    "nsp": project.Nsp,
    "config": project.ProjectConfig.AsJSON(),
  }
}

func (project *Project) GetHost() projecthost.Host {
  var host projecthost.Host

  switch project.Host {
  case projecthost.GitHubName:
    return &projecthost.GitHub{Project: project}
  }

  host.Configure()
  return host
}

func (project *Project) SplitNsp() []string {
  return strings.Split(project.Nsp, "/")
}

func (project *Project) Owner() string {
  return project.SplitNsp()[1]
}

func (project *Project) Repo() string {
  return project.SplitNsp()[2]
}

func hostNameForNsp(nsp string) string {
  host := ""

  switch true {
  case strings.HasPrefix(nsp, projecthost.GitHubDomain):
    host = projecthost.GitHubName
  }

  return host
}