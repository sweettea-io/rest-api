package model

import (
  "fmt"
  "strings"
  "github.com/jinzhu/gorm"
  "github.com/sweettea-io/rest-api/internal/pkg/util/enc"
  "github.com/sweettea-io/rest-api/internal/pkg/util/unique"
)

/*
  has_one -- ProjectConfig

  has_many --> Commits
  has_many --> Models
*/
type Project struct {
  gorm.Model
  Uid             string        `gorm:"type:varchar(240);default:null;not null;unique;index:project_uid"`
  Nsp             string        `gorm:"type:varchar(360);default:null;not null;unique;index:project_nsp"`
  ProjectConfig   *ProjectConfig
  ProjectConfigID uint          `gorm:"default:null;not null;index:project_config_id"`
  Commits         []Commit
  Models          []Model
}

// Downcase project namespace before saving.
func (project *Project) BeforeSave(scope *gorm.Scope) error {
  scope.SetColumn("Nsp", strings.ToLower(project.Nsp))
  return nil
}

// Assign Uid to Project before creation.
func (project *Project) BeforeCreate(scope *gorm.Scope) error {
  scope.SetColumn("Uid", unique.NewUid())
  return nil
}

func (project *Project) AsJSON() enc.JSON {
  return enc.JSON{
    "nsp": project.Nsp,
    "createdAt": project.CreatedAt,
  }
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

func (project *Project) Url() string {
  return fmt.Sprintf("https://%s", project.Nsp)
}