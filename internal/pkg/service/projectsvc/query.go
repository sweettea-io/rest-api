package projectsvc

import (
  "fmt"
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/pkg/model"
)

// All returns all Project records ordered by nsp
func All() []model.Project {
  projects := []model.Project{}
  app.DB.Order("nsp desc").Find(&projects)
  return projects
}

// FromNsp attempts to find a Project record for the given nsp.
// Will return an error if no record is found.
func FromNsp(nsp string) (*model.Project, error) {
  // Find Project by nsp.
  var project model.Project
  result := app.DB.Where(&model.Project{Nsp: nsp}).Find(&project)

  // Return error if not found.
  if result.RecordNotFound() {
    return nil, fmt.Errorf("Project(nsp=%s) not found.\n", nsp)
  }

  return &project, nil
}